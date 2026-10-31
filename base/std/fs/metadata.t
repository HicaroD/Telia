package fs

enum FileType u32 {
    Regular = 1
    Directory = 2
    Symlink = 3
    Other = 4
}

struct Metadata {
    kind FileType
    volume u64
    identity u64
    size u64
    modified_seconds i64
    modified_nanoseconds u32
    permissions u32
}

fn stat(path string) (Metadata, error) {
    return metadata(path, true)
}

fn lstat(path string) (Metadata, error) {
    return metadata(path, false)
}

fn metadata(path string, follow bool) (Metadata, error) {
    result := Metadata{}
    err := validate_path(path)
    if err != nil {
        return result, err
    }
    kind u32 := 0
    err = metadata_into(path, follow, &kind, &result.volume, &result.identity, &result.size, &result.modified_seconds, &result.modified_nanoseconds, &result.permissions)
    result.kind = FileType::Other
    if kind == 1 {
        result.kind = FileType::Regular
    }
    if kind == 2 {
        result.kind = FileType::Directory
    }
    if kind == 3 {
        result.kind = FileType::Symlink
    }
    return result, err
}

fn same_file(left Metadata, right Metadata) bool {
    return left.volume == right.volume and left.identity == right.identity
}

@[default_cc="c"]
extern {
    @[link_name="_telia_fs_metadata"]
    fn metadata_into(path string, follow bool, kind *u32, volume *u64, identity *u64, size *u64, seconds *i64, nanos *u32, permissions *u32) error
}
