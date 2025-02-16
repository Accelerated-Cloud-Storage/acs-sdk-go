package fuse

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

    clientPkg "github.com/AcceleratedCloudStorage/acs-sdk-go/client"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
	"golang.org/x/sys/unix"
)

// Regression test: create/delete/list files, create/delete/list folder, create folder with files and other folders and then list and delete them, check path updates as you cd in and out folders, check files contents, move/copy files in and out of folders
type FS struct {
    client     *clientPkg.ACSClient
    bucketName string
}

func (f *FS) Root() (fs.Node, error) {
    return &Dir{
        fs: f,
        path: "",
    }, nil
}

type Dir struct {
    fs   *FS
    path string
}

type File struct {
    fs     *FS
    name   string
    size   uint64
    mtime  time.Time
}

type FileHandle struct {
    file  *File
}

var (
    mountUID int
    mountGID int
    globalLock sync.Mutex
    
    // Add caching maps
    headCache    = make(map[string]*headCacheEntry)
    headCacheMu  sync.RWMutex
)

type headCacheEntry struct {
    exists    bool
    isDir     bool
    timestamp time.Time
}

// Add cache helper functions
func getCachedHead(path string) (exists bool, isDir bool, found bool) {
    headCacheMu.RLock()
    defer headCacheMu.RUnlock()
    
    if entry, ok := headCache[path]; ok {
        // Cache entries expire after 60 seconds
        if time.Since(entry.timestamp) < 60*time.Second {
            return entry.exists, entry.isDir, true
        }
    }
    return false, false, false
}

func setCachedHead(path string, exists bool, isDir bool) {
    headCacheMu.Lock()
    defer headCacheMu.Unlock()
    
    headCache[path] = &headCacheEntry{
        exists:    exists,
        isDir:     isDir,
        timestamp: time.Now(),
    }
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return syscall.ENOENT
    }
    log.Printf("Dir.Attr called for path: %s", d.path)
    a.Inode = 1
    a.Mode = os.ModeDir | 0755
    a.Uid = uint32(mountUID)
    a.Gid = uint32(mountGID)
    a.Size = 0
    // Set reasonable directory timestamps
    now := time.Now()
    a.Atime = now
    a.Mtime = now
    a.Ctime = now
    return nil
}

// Optimized Lookup implementation
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return nil, syscall.ENOENT
    }
    log.Printf("lookup: Searching for file %s", name)
    // Ignore system files
    if name == "Input" || name == ".Trash" || name == ".Trash-1000" || 
       name == ".xdg-volume-info" || name == "autorun.inf" {
        return nil, syscall.ENOENT
    }
    
    // Construct full path
    fullPath := name
    if d.path != "" {
        fullPath = d.path + "/" + name
    }
    
    // Check cache first
    if exists, isDir, found := getCachedHead(fullPath); found {
        if (!exists) {
            return nil, syscall.ENOENT
        }
        if (isDir) {
            return &Dir{fs: d.fs, path: fullPath}, nil
        }
        return &File{fs: d.fs, name: fullPath}, nil
    }

    log.Printf("lookup: Checking for file %s in bucket %s", name, d.fs.bucketName)
    
    // First check if it's a file (most common case)
    if _, err := d.fs.client.HeadObject(ctx, d.fs.bucketName, fullPath); err == nil {
        setCachedHead(fullPath, true, false)
        return &File{fs: d.fs, name: fullPath}, nil
    }

    log.Printf("lookup: Checking for directory %s in bucket %s", name, d.fs.bucketName)
    
    // Then check for directory with trailing slash
    dirPath := fullPath + "/"
    if _, err := d.fs.client.HeadObject(ctx, d.fs.bucketName, dirPath); err == nil {
        setCachedHead(fullPath, true, true)
        return &Dir{fs: d.fs, path: fullPath}, nil
    }

    // If not found as file or explicit directory, do a single list operation
    // with a limit of 1 to check for prefix matches
    objects, listErr := d.fs.client.ListObjects(ctx, d.fs.bucketName, &clientPkg.ListObjectsOptions{
        Prefix:    fullPath + "/",
        MaxKeys:   1,
    })
    if listErr == nil && len(objects) > 0 {
        setCachedHead(fullPath, true, true)
        return &Dir{fs: d.fs, path: fullPath}, nil
    }
    
    setCachedHead(fullPath, false, false)
    log.Printf("lookup: Searched for file %s", name)
    return nil, syscall.ENOENT
}

// Add Open method to Dir to support directory operations
func (d *Dir) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return nil, syscall.ENOENT
    }
    log.Printf("Dir.Open: %s", d.path)
    return d, nil
}

// Implement Handle interface for Dir
func (d *Dir) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return syscall.ENOENT
    }
    log.Printf("Dir.Release: %s", d.path)
    return nil
}

func (d *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return nil, syscall.ENOENT
    }
    
    // Construct full path
    fullPath := req.Name
    if d.path != "" {
        fullPath = d.path + "/" + req.Name
    }
    dirPath := fullPath + "/"
    
    log.Printf("Mkdir: Creating directory %s in bucket %s", dirPath, d.fs.bucketName)
    
    err := d.fs.client.PutObject(ctx, d.fs.bucketName, dirPath, []byte{})
    if err != nil {
        log.Printf("Mkdir failed for %s: %v", dirPath, err)
        return nil, err
    }

    log.Printf("Mkdir: Created directory %s in bucket %s", dirPath, d.fs.bucketName)
    
    // Return the new directory node with its path
    return &Dir{fs: d.fs, path: fullPath}, nil
}

func (d *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return syscall.ENOENT
    }
    // Construct full path
    fullPath := req.Name
    if d.path != "" {
        fullPath = d.path + "/" + req.Name
    }

    // Add trailing slash for directories
    if req.Dir {
        fullPath += "/"
    }

    log.Printf("Remove: Deleting %s (dir: %v) from bucket %s", fullPath, req.Dir, d.fs.bucketName)
    
    // First verify the object exists
    _, err := d.fs.client.HeadObject(ctx, d.fs.bucketName, fullPath)
    if err != nil {
        log.Printf("Remove: Object not found: %s, error: %v", fullPath, err)
        return syscall.ENOENT
    }

    log.Printf("Remove: Deleting %s (dir: %v) from bucket %s and object exists", fullPath, req.Dir, d.fs.bucketName)
    
    // If it's a directory, make sure it's empty
    if req.Dir {
        objects, err := d.fs.client.ListObjects(ctx, d.fs.bucketName, &clientPkg.ListObjectsOptions{
            Prefix: fullPath,
        })
        if err != nil {
            log.Printf("Remove: Failed to list directory contents: %v", err)
            return err
        }
        
        // Check if directory has any contents (excluding the directory marker itself)
        for _, obj := range objects {
            if obj != fullPath && strings.HasPrefix(obj, fullPath) {
                log.Printf("Remove: Directory not empty: %s", fullPath)
                return syscall.ENOTEMPTY
            }
        }
    }

    log.Printf("Remove: Deleting %s (dir: %v) from bucket %s and directory is empty", fullPath, req.Dir, d.fs.bucketName)
    
    // Perform the deletion
    err = d.fs.client.DeleteObject(ctx, d.fs.bucketName, fullPath)
    if err != nil {
        log.Printf("Remove: Delete failed for %s: %v", fullPath, err)
        return err
    }
    
    log.Printf("Remove: Successfully deleted %s", fullPath)
    return nil
}

func (d *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (node fs.Node, handle fs.Handle, err error) {    
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return nil, nil, syscall.ENOENT
    }

    fullPath := req.Name
    if d.path != "" {
        fullPath = d.path + "/" + req.Name
    }
    log.Printf("Create: Starting create for file %s", fullPath)

    // Create empty file in S3
    if err := d.fs.client.PutObject(ctx, d.fs.bucketName, fullPath, []byte{}); err != nil {
        log.Printf("Create failed: %v", err)
        return nil, nil, err
    }

    f := &File{
        fs:     d.fs,
        name:   fullPath,
        size:   0,
        mtime:  time.Now(),
    }

    resp.Attr = fuse.Attr{
        Inode: 2,
        Mode:  req.Mode,
        Size:  f.size,
        Mtime: f.mtime,
        Uid:   uint32(mountUID),
        Gid:   uint32(mountGID),
    }
    
    handle = &FileHandle{
        file:  f,
    }

    log.Printf("Create: Successfully created file %s", fullPath)
    return f, handle, nil
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
    globalLock.Lock()
    defer globalLock.Unlock()

    if f == nil {
        return syscall.ENOENT
    }

    log.Printf("Attr: Starting attribute retrieval for %s", f.name)

    a.Inode = 2
    a.Mode = 0666
    a.Uid = uint32(mountUID)
    a.Gid = uint32(mountGID)
    a.Size = f.size
    a.Blocks = (f.size + 511) / 512
    a.Atime = time.Now()
    a.Mtime = f.mtime
    a.Ctime = f.mtime

    log.Printf("Attr: Successfully returned attributes for %s", f.name)
    return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
    globalLock.Lock()
    defer globalLock.Unlock()
    if f == nil {
        return nil, syscall.ENOENT
    }
    log.Printf("ReadAll: Reading all bytes from %s", f.name)
    // Use the correct filename instead of hardcoded "file-name"
    data, err := f.fs.client.GetObject(ctx, f.fs.bucketName, f.name)
    if (err != nil) {
        return nil, err
    }
    log.Printf("ReadAll: Read all bytes from %s", f.name)
    return data, nil
}

// Implement Read method for File
func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if f == nil {
        return syscall.ENOENT
    }
    log.Printf("Read: Reading bytes from file %s", f.name)
    // Get the full file data
    data, err := f.fs.client.GetObject(ctx, f.fs.bucketName, f.name)
    if err != nil {
        log.Printf("Read error for %s: %v", f.name, err)
        return err
    }

    // Validate offset
    if req.Offset < 0 || req.Offset >= int64(len(data)) {
        resp.Data = []byte{}
        return nil
    }

    // Calculate end position
    end := req.Offset + int64(req.Size)
    if end > int64(len(data)) {
        end = int64(len(data))
    }

    // Return the requested portion of data
    resp.Data = data[req.Offset:end]

    log.Printf("Read: Read bytes from file %s", f.name)
    return nil
}

// Open method for File to support read/write operations
func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {   
    globalLock.Lock()
    defer globalLock.Unlock() 
    // Add error checking
    if f == nil {
        return nil, syscall.ENOENT
    }
    log.Printf("File.Open: Opening File %s", f.name)   
    
    // Create new handle with proper initialization
    handle := &FileHandle{
        file: f,
    }

    // Verify the file exists in bucket
    _, err := f.fs.client.HeadObject(ctx, f.fs.bucketName, f.name)
    if err != nil {
        log.Printf("File.Open: Failed to verify file existence: %v", err)
        return nil, syscall.ENOENT
    }

    log.Printf("File.Open: Successfully opened %s", f.name)
    return handle, nil
}

// Implement Handle interface for File
func (f *File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if f == nil {
        return syscall.ENOENT
    }
    log.Printf("File.release: %s", f.name)
    return nil
}

// Add after other File methods
func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if f == nil {
        return syscall.ENOENT
    }
    log.Printf("File.Flush: %s", f.name)
    return nil
}

// Write method for FileHandle
func (fh *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if fh == nil {
        return syscall.ENOENT
    }
    log.Printf("Write: Writing %d bytes to %s", len(req.Data), fh.file.name)

    // Sync write to S3
    err := fh.file.fs.client.PutObject(ctx, fh.file.fs.bucketName, fh.file.name, req.Data)
    if err != nil {
        log.Printf("Write failed for %s: %v", fh.file.name, err)
    }

    
    fh.file.size = uint64(len(req.Data))
    fh.file.mtime = time.Now()
    resp.Size = len(req.Data)

    log.Printf("Write: Wrote %d bytes to %s to S3", len(req.Data), fh.file.name)
    return nil
}

func (fh *FileHandle) Open(ctx context.Context, req *fuse.ReleaseRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if fh == nil {
        return syscall.ENOENT
    }
    log.Printf("Open: Opened file %s", fh.file.name)
    return nil
}

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if fh == nil {
        return syscall.ENOENT
    }
    log.Printf("Release: Closed file %s", fh.file.name)
    return nil
}

// Read method for FileHandle
func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if fh == nil {
        return syscall.ENOENT
    }
    log.Printf("Read: Reading %d bytes from file handle %s from bucket %s", req.Size, fh.file.name, fh.file.fs.bucketName)

    // Get the data - returns error if file not found
    data, err := fh.file.fs.client.GetObject(ctx, fh.file.fs.bucketName, fh.file.name)
    if err != nil {
        log.Printf("Read error for %s: %v", fh.file.name, err)
        return err 
    }

    // Check bounds
    if req.Offset < 0 || req.Offset >= int64(len(data)) || req.Size < 0 {
        resp.Data = []byte{}
        return nil
    }

    // Calculate safe read bounds
    end := req.Offset + int64(req.Size)
    if end > int64(len(data)) {
        end = int64(len(data))
    }

    // Copy data to response
    resp.Data = make([]byte, end-req.Offset)
    copy(resp.Data, data[req.Offset:end])

    log.Printf("Read: Read %d bytes from file handle %s from S3", req.Size, fh.file.name)
    return nil
}

// Add after other FileHandle methods
func (fh *FileHandle) Flush(ctx context.Context, req *fuse.FlushRequest) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if fh == nil {
        return syscall.ENOENT
    }
    log.Printf("FileHandle.Flush: %s", fh.file.name)
    return nil
}

// Get all foles in directory 
func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return nil, syscall.ENOENT
    }
    log.Printf("ReadDirAll: Listing objects in bucket %s with prefix %s", d.fs.bucketName, d.path)
    
    prefix := ""
    if d.path != "" {
        prefix = d.path + "/"
    }
    
    // List all objects with the prefix
    objects, err := d.fs.client.ListObjects(ctx, d.fs.bucketName, &clientPkg.ListObjectsOptions{
        Prefix: prefix,
    })
    if err != nil {
        log.Printf("ReadDirAll: ListObjects error: %v", err)
        return nil, err
    }

    dirents := []fuse.Dirent{{
        Inode: 1,
        Type:  fuse.DT_Dir,
        Name:  ".",
    }}

    seen := make(map[string]bool)
    i := uint64(2)

    // Process all objects
    for _, obj := range objects {
        // Skip empty objects and the prefix itself
        if (obj == "" || obj == prefix) {
            continue
        }

        // Get the relative path by removing the prefix
        name := obj
        if prefix != "" {
            if len(obj) <= len(prefix) {
                continue
            }
            name = obj[len(prefix):]
        }

        // Split into components to handle nested paths
        components := strings.Split(name, "/")
        if len(components) == 0 {
            continue
        }

        // Get the first component
        firstComp := components[0]
        if firstComp == "" {
            continue
        }

        // Skip if already processed
        if seen[firstComp] {
            continue
        }
        seen[firstComp] = true

        // Determine if it's a directory (has more components or ends with slash)
        isDir := len(components) > 1 || strings.HasSuffix(obj, "/")

        log.Printf("Adding entry: %s (isDir: %v)", firstComp, isDir)
        dirents = append(dirents, fuse.Dirent{
            Inode: i,
            Type:  map[bool]fuse.DirentType{true: fuse.DT_Dir, false: fuse.DT_File}[isDir],
            Name:  firstComp,
        })
        i++
    }

    log.Printf("ReadDirAll: Returning %d entries for %s", len(dirents), d.path)
    return dirents, nil
}

// Rename a file similar to the MV comamnd 
func (d *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
    globalLock.Lock()
    defer globalLock.Unlock()
    if d == nil {
        return syscall.ENOENT
    }
    // Get old and new paths
    oldPath := req.OldName
    if d.path != "" {
        oldPath = d.path + "/" + req.OldName
    }

    // Get the new directory path from the target Dir node
    newDirPath := ""
    if newDir, ok := newDir.(*Dir); ok {
        newDirPath = newDir.path
    }
    // Construct new path
    newPath := req.NewName
    if newDirPath != "" {
        newPath = newDirPath + "/" + req.NewName
    }

    log.Printf("Rename: Moving %s to %s", oldPath, newPath)

    // Check if source exists
    data, err := d.fs.client.GetObject(ctx, d.fs.bucketName, oldPath)
    if err != nil {
        // Check if it's a directory
        dirPath := oldPath + "/"
        data, err = d.fs.client.GetObject(ctx, d.fs.bucketName, dirPath)
        if err != nil {
            log.Printf("Rename: Source not found: %s", oldPath)
            return syscall.ENOENT
        }
        // It's a directory, append slash to new path
        newPath += "/"
    }

    // Create the new object
    err = d.fs.client.PutObject(ctx, d.fs.bucketName, newPath, data)
    if err != nil {
        log.Printf("Rename: Failed to create new object: %v", err)
        return err
    }

    // Delete the old object
    err = d.fs.client.DeleteObject(ctx, d.fs.bucketName, oldPath)
    if err != nil {
        log.Printf("Rename: Failed to delete old object: %v", err)
        return err
    }

    return nil
}

func NewFUSEMount(bucketname string, mountPoint string) error {
    var err error
    client, err := clientPkg.NewClient()
    if err != nil {
        return err
    }

    // Ensure mount point exists with correct permissions
    fmt.Println("Mounting bucket", bucketname, "at", mountPoint)
    if err := os.MkdirAll(mountPoint, 0777); err != nil {
        return fmt.Errorf("failed to create mount point: %v", err)
    }
        
    if _, err := os.Stat(mountPoint); os.IsNotExist(err) {
        log.Fatalf("Mount point %s does not exist", mountPoint)
    }
    
    if err := unix.Access(mountPoint, unix.W_OK); err != nil {
        log.Fatalf("Mount point %s is not writable: %v", mountPoint, err)
    }

    // Get current user ID and group ID
    mountUID = os.Getuid()
    mountGID = os.Getgid()

    // Set correct ownership
    if err := os.Chown(mountPoint, mountUID, mountGID); err != nil {
        return fmt.Errorf("failed to change mount point ownership: %v", err)
    }

    // Mount with minimal supported options
    c, err := fuse.Mount(
        mountPoint,
        fuse.FSName("objectstoragecache"),
        fuse.Subtype("objectstoragecachefs"),
        fuse.DefaultPermissions(),
        fuse.MaxBackground(10000),
        fuse.CongestionThreshold(7000),
    )
    if err != nil {
        log.Printf("Failed to mount: %v", err)
        return err
    }
    defer c.Close()

    filesys := &FS{
        client: client,
        bucketName: bucketname,
    }
    if err := fs.Serve(c, filesys); err != nil {
        log.Printf("Failed to serve: %v", err)
        fuse.Unmount(mountPoint)
        return err
    }

    return nil
}
