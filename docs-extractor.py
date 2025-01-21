import argparse
import textwrap
from os import path
import glob
import os
import re
from dataclasses import dataclass


tabs = re.compile(r'<Tabs groupId="language">(.*?)</Tabs>', re.DOTALL)
tabitems = re.compile(r'<TabItem value="(.*?)">(.*?)</TabItem>', re.DOTALL)
snippet = re.compile(r"```(\S+).*?(?: file=(\S*))?")
codeblock = re.compile(r"```[^\n]*(.*?)```", re.DOTALL)
headers = re.compile(r"^(#+) (.*)", re.MULTILINE)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("output")
    args = parser.parse_args()

    examples = []
    for filename in glob.glob("docs/**/*.mdx", recursive=True):
        with open(filename) as docsfile:
            data = docsfile.read()
            clean_data = data

            # HACK: trim out the codeblock data (but store the original data)
            # this ensures we don't interpret code block contents as markdown data
            blocks = {}
            for blockmatch in codeblock.finditer(clean_data):
                blocks[blockmatch.start()] = blockmatch.group(1)
                clean_data = (
                    clean_data[: blockmatch.start(1)]
                    + re.sub(r"\S", " ", blockmatch.group(1))
                    + clean_data[blockmatch.end(1) :]
                )

            lasttabsmatch = None
            for tabsmatch in tabs.finditer(clean_data):
                files = []
                for tabmatch in tabitems.finditer(
                    clean_data, tabsmatch.start(1), tabsmatch.end(1)
                ):
                    snippetmatch = snippet.search(
                        clean_data, tabmatch.start(2), tabmatch.end(2)
                    )
                    if not snippetmatch:
                        continue
                    language, snippath = snippetmatch.group(1), snippetmatch.group(2)
                    if snippath:
                        snippath = path.join(path.dirname(filename), snippath)
                        files.append(SourceFile(language, filename=snippath))
                    else:
                        files.append(
                            SourceFile(language, raw=blocks[snippetmatch.start()])
                        )

                hs = list(
                    headers.finditer(
                        clean_data[: tabsmatch.endpos], endpos=tabsmatch.start()
                    )
                )
                relevant = []
                for h in hs:
                    level, title = len(h.group(1)) - 1, h.group(2)
                    relevant = relevant[:level]
                    relevant.append(title)

                description = data[
                    max(
                        hs[-1].end() if hs else 0,
                        lasttabsmatch.end() if lasttabsmatch else 0,
                    ) : tabsmatch.start()
                ]
                description = description.strip()

                lasttabsmatch = tabsmatch

                example = Example(relevant, description, files)
                langs = {f.language for f in example.files}
                if len(langs) < 2:  # must have at least two different languages
                    continue
                if "go" not in langs:  # must contain go
                    continue
                examples.append(example)

    os.makedirs(args.output, exist_ok=True)

    slugs = set()
    for example in examples:
        baseslug = example.slug
        slug = baseslug
        n = 1
        while slug in slugs:
            n += 1
            slug = baseslug + "-" + str(n)
        slugs.add(slug)

        # with open(path.join(args.output, slug) + ".md", "w") as f:
        #     f.write("# " + example.headers[-1] + "\n\n")
        #     f.write(example.content + "\n\n")
        #
        #     for source in example.files:
        #         f.write("```" + source.language + "\n")
        #         f.write(source.content)
        #         f.write("```\n\n")

        with open(path.join(args.output, slug) + ".xml", "w") as f:
            f.write("<example>\n")
            f.write("  <name>" + example.headers[-1] + "</name>\n")
            f.write(
                "  <description>\n"
                + textwrap.indent(example.content, "    ")
                + "\n  </description>\n"
            )

            for source in example.files:
                f.write(f'  <snippet language="{source.language}">\n')
                f.write(textwrap.indent(source.content, "    "))
                f.write("  </snippet>\n")

            f.write("</example>\n")


@dataclass
class SourceFile:
    language: str
    filename: str | None = None
    raw: str | None = None

    @property
    def content(self) -> str:
        if self.raw:
            return self.raw

        if self.filename:
            with open(self.filename) as f:
                return f.read()

        raise ValueError("no filename or contents")


@dataclass
class Example:
    headers: list[str]
    content: str
    files: list[SourceFile]

    @property
    def slug(self):
        return (
            re.sub(r"([a-zA-Z0-9]+)[^a-zA-Z0-9]+", r"\1-", self.headers[-1])
            .strip("-")
            .lower()
        )


if __name__ == "__main__":
    main()
