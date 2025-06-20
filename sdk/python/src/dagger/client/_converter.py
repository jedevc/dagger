import enum

from cattrs.preconf.json import make_converter as make_json_converter


def make_converter():
    conv = make_json_converter(
        detailed_validation=True,
    )
    register_enum_hooks(conv)
    return conv


def register_enum_hooks(conv):
    conv.register_unstructure_hook_func(_is_enum, _to_enum_name)
    conv.register_structure_hook_func(_is_enum, _from_enum_name)


def _to_enum_name(val: enum.Enum) -> str:
    return val.name


def _from_enum_name(name: str, cls: type[enum.Enum]) -> enum.Enum:
    return cls[name]


def _is_enum(t) -> bool:
    return issubclass(t, enum.Enum)
