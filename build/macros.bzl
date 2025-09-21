
def lit2md(name, src, out):
    native.genrule(
        name = name,
        srcs = [src],
        outs = [out],
        tools = [
            Label("@lit2md//cmd/lit2md"),
        ],
        cmd = """$(location @lit2md//cmd/lit2md) \
                --input=$(location {src}) \
                --output=$(location {out})""".format(src=src,out=out),
    )


def lit2md2(name, src, out):
    lit2md(name=name, src=src, out="x.{out}".format(out=out))

