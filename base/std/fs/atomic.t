package fs

use "std::strings"
use "std::path"
use "std::errors"

fn create_temp(destination string) (File, strings::Owned, error) {
    file := File{}
    name := strings::Owned{}
    err := validate_path(destination)
    if err == nil {
        err = temp_into(destination, &file.state, &name.text)
    }
    return file, name, err
}

fn write(file *File, data string) error {
    return file_write(file.state, data)
}

fn sync(file *File) error {
    return file_sync(file.state)
}

fn preserve_permissions(source string, destination *File) error {
    err := validate_path(source)
    if err != nil {
        return err
    }
    return permissions_into(source, destination.state)
}

fn rename(source string, destination string) error {
    err := validate_path(source)
    if err == nil {
        err = validate_path(destination)
    }
    if err != nil {
        return err
    }
    return rename_native(source, destination)
}

fn remove(path string) error {
    err := validate_path(path)
    if err != nil {
        return err
    }
    return remove_native(path)
}

fn sync_parent(path string) error {
    err := validate_path(path)
    if err != nil {
        return err
    }
    return parent_sync(path)
}

fn cleanup_temp(name *strings::Owned) {
    if strings::len(strings::view(name)) != 0 {
        // Cleanup must still work after the allocator has failed.
        cleanup_native(strings::view(name))
    }
}

// The boolean reports replacement even when subsequent directory sync fails.
fn write_atomic(destination string, data string) (bool, error) {
    target, target_err := path::canonicalize_for_create(destination)
    if target_err != nil {
        return false, target_err
    }
    defer strings::free(&target)
    target_path := strings::view(&target)
    info, info_err := stat(target_path)
    if info_err != nil and info_err.kind != errors::Error::NotFound {
        return false, info_err
    }
    if info_err == nil and info.kind != FileType::Regular {
        return false, error(errors::Error::InvalidPath, "atomic save requires a regular file")
    }
    file, temporary, err := create_temp(target_path)
    if err != nil {
        return false, err
    }
    defer strings::free(&temporary)
    temporary_path := strings::view(&temporary)
    defer cleanup_temp(&temporary)
    defer close(&file)
    err = write(&file, data)
    if err == nil and info_err == nil {
        err = preserve_permissions(target_path, &file)
        // Windows replacement preserves the target DACL through ReplaceFileW.
        if err.kind == errors::Error::Unsupported {
            err = zero<error>()
        }
    }
    if err == nil {
        err = sync(&file)
    }
    close_err := close(&file)
    if err == nil {
        err = close_err
    }
    if err != nil {
        return false, err
    }
    err = rename(temporary_path, target_path)
    if err != nil {
        return false, err
    }
    strings::free(&temporary)
    return true, sync_parent(target_path)
}

@[default_cc="c"]
extern {
    @[link_name="_telia_file_temp"]
    fn temp_into(destination string, state *rawptr, name *string) error
    @[link_name="_telia_file_write"]
    fn file_write(state rawptr, data string) error
    @[link_name="_telia_file_sync"]
    fn file_sync(state rawptr) error
    @[link_name="_telia_file_permissions"]
    fn permissions_into(source string, state rawptr) error
    @[link_name="_telia_file_rename"]
    fn rename_native(source string, destination string) error
    @[link_name="_telia_file_remove"]
    fn remove_native(path string) error
    @[link_name="_telia_parent_sync"]
    fn parent_sync(path string) error
    @[link_name="_telia_file_cleanup_temp"]
    fn cleanup_native(path string) error
}
