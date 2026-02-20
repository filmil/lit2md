load("//build:macros.bzl", _lit2md = "lit2md", _lit2md2 = "lit2md2")
load("//build:my_package_name.bzl", _name_part_from_command_line = "name_part_from_command_line")

lit2md = _lit2md
lit2md2 = _lit2md2
name_part_from_command_line = _name_part_from_command_line

