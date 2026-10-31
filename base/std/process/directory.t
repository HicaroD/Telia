package process

use "std::path"

fn chdir(directory string) error {
    err := path::validate(directory)
    if err != nil {
        return err
    }
    return chdir_native(directory)
}

@[default_cc="c"]
extern {
    @[link_name="_telia_chdir"]
    fn chdir_native(directory string) error
}
