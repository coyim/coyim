#!/usr/bin/env bash

new_pc_files_dir=$1
mkdir -p "$new_pc_files_dir"

fix_pc_file()
{
    pc_file_destination=$1
    pc_file=$2

    pc_file_name="$(basename $pc_file)"
    sed -e "s/prefix=$prefix/prefix=$replace/g" -e 's/-Wl,-luuid/-luuid/g' -e '$aLDFLAGS: -Wl' "$pc_file" > "$pc_file_destination/$pc_file_name"
}

pkg_config_paths="$(pkg-config --variable pc_path pkg-config)"
IFS=';' read -ra pkg_config_paths_arr <<< "$pkg_config_paths"
pc_files="$(find "${pkg_config_paths_arr[@]}" \( -name 'gdk-2.0.pc' -o -name 'gdk-win32-2.0.pc' -o -name 'gdk-3.0.pc' -o -name 'gdk-win32-3.0.pc' \))"

while IFS= read -r pc_file; do
    echo "Fixing \"$pc_file\" and copying it to the \"$new_pc_files_dir\" folder"
    fix_pc_file "$new_pc_files_dir" "$pc_file"
done <<< "$pc_files"

