
def lit2md(name, src, out):
    native.genrule(
        name = name,
        srcs = [src],
        outs = [out],
        tools = [
            Label("//cmd/lit2md"),
        ],
        cmd = """$(location //cmd/lit2md) \
                --input=$(location {src}) \
                --output=$(location {out})""".format(src=src,out=out),
    )
