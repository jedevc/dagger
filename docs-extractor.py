import argparse
from os import path
import glob
import os
import re
from dataclasses import dataclass


tabs = re.compile(r'<Tabs groupId="language">(.*?)</Tabs>', re.DOTALL)
tabitems = re.compile(r'<TabItem value="(\S*?)">(.*?)</TabItem>', re.DOTALL)
snippet = re.compile(r"```(\S+).* file=(\S*)")
codeblock = re.compile(r"```[^\n]*(.*?)```", re.DOTALL)
headers = re.compile(r"^(#+) (.*)", re.MULTILINE)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("output")
    args = parser.parse_args()

    examples = []
    for filename in glob.glob("docs/**/*.mdx", recursive=True):
        if "cookbook" not in filename:
            continue

        with open(filename) as docsfile:
            data = docsfile.read()

            # HACK: trim out the codeblock data
            for blockmatch in codeblock.finditer(data):
                data = (
                    data[: blockmatch.start(1)]
                    + re.sub(r"\S", " ", blockmatch.group(1))
                    + data[blockmatch.end(1) :]
                )

            lasttabsmatch = None
            for tabsmatch in tabs.finditer(data):
                files = []
                for tabmatch in tabitems.finditer(tabsmatch.group(1)):
                    snippetmatch = snippet.search(tabmatch.group(2))
                    if not snippetmatch:
                        continue
                    language, snippath = snippetmatch.group(1), snippetmatch.group(2)
                    snippath = path.join(path.dirname(filename), snippath)
                    files.append(SourceFile(language, snippath))

                hs = list(
                    headers.finditer(data[: tabsmatch.endpos], endpos=tabsmatch.start())
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

                examples.append(Example(relevant, description, files))

                lasttabsmatch = tabsmatch

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

        with open(path.join(args.output, slug) + ".md", "w") as f:
            f.write("# " + example.headers[-1] + "\n\n")
            f.write(example.content + "\n\n")

            for source in example.files:
                f.write("```" + source.language + "\n")
                f.write(source.contents)
                f.write("```\n\n")


@dataclass
class SourceFile:
    language: str
    filename: str

    # TODO: contents might actually be inline

    @property
    def contents(self):
        with open(self.filename) as f:
            return f.read()


@dataclass
class Example:
    headers: list[str]
    content: str
    files: list[SourceFile]

    @property
    def slug(self):
        return (
            re.sub(r"([a-zA-Z0-9]+)[^a-zA-Z0-9]*", r"\1-", self.headers[-1])
            .strip("-")
            .lower()
        )


if __name__ == "__main__":
    main()
