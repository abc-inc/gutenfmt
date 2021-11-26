<h1 align="center">gutenfmt</h1>

<p align="center">
  <i>gutenfmt is a cross-platform command line tool that converts name-value pairs and JSON to different output formats,
    <br>including, but not limited to, CSV, JSON, YAML or a table with aligned columns.</i>
  <br>
</p>

<p align="center">
  <a href="https://abc-inc.github.io/gutenfmt/" target="_blank"><strong>abc-inc.github.io/gutenfmt</strong></a>
  <br>
</p>

## Examples

The following examples provide a brief overview of gutenfmt and its features.

### List OpenJDK Packages on a Debian-System as Colorized JSON

```shell
$ dpkg-query -W -f='${Package}\t${Version}\t${Status}\n' "openjdk*" | grep -Fv unknown | cut -f 1,2 | gutenfmt

{
  "openjdk-16-jdk": "16.0.2+7-2",
  "openjdk-16-jdk-headless": "16.0.2+7-2",
  "openjdk-16-jre": "16.0.2+7-2",
  "openjdk-16-jre-headless": "16.0.2+7-2",
  "openjdk-16-source": "16.0.2+7-2",
  "openjdk-8-jdk": "8u302-b08-0ubuntu2",
  "openjdk-8-jdk-headless": "8u302-b08-0ubuntu2",
  "openjdk-8-jre": "8u302-b08-0ubuntu2",
  "openjdk-8-jre-headless": "8u302-b08-0ubuntu2"
}
```

### Print JavaScript Dependencies as Table With Aligned Columns

```shell
$ jq '.dependencies' package.json | tr -d '" ,' | gutenfmt -o text | sort

...
@types/yargs                         ^17.0.3
@webcomponents/custom-elements       ^1.1.0
xhr2                                 0.2.1
yargs                                ^17.2.1
```

## Further Information

gutenfmt is a also a Go library that helps to create applications, which support different output formats.
It comes with an opinionated view of the third-party libraries, so that you can get started with minimum fuss.
In other words, Formatters come with meaningful defaults.

The primary goals are:

- Provide a radically faster and widely accessible getting-started experience for Command Line Interface (CLI) development.
- Be opinionated out of the box but get out of the way quickly as requirements start to diverge from the defaults.
- Provide a range of standardized output formats that are common to state-of-the-art projects.

Both, User Manual and Development Guide, can be found at
<a href="https://abc-inc.github.io/gutenfmt/" target="_blank"><strong>abc-inc.github.io/gutenfmt</strong></a>.
