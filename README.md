# scriptura

A CLI app for reading the Bible by book, chapter(s), or verse(s); uses the Douay-Rheims 1899 American Edition (DRA) version of the Bible.

## Usage

This usage information can be found by running `scriptura --help`:

```
Usage: scriptura <book>
       scriptura <book> [start_chapter]-[end_chapter]
       scriptura <book> <chapter>
       scriptura <book> <chapter>:[start_verse]-[end_verse]
       scriptura <book> <chapter>:<verse>

Options:
       --books    Print the list of books and exit
       --help     Print this message and exit
       --license  Print license, citation information and exit
       --version  Print version and exit
```

### Usage examples

<pre>
<code>>>> scriptura Genesis 12:1
And the Lord said to Abram: Go forth out of thy country, and from thy kindred, and out of thy father’s house, and come into the land which I shall shew thee.
</code></pre>
<pre>
<code>>>> scriptura Acts 5:40-42
  <b>40</b> And calling in the apostles, after they had scourged them, they charged them that they should not speak at all in the name of Jesus; and they dismissed them.
  <b>41</b> And they indeed went from the presence of the council, rejoicing that they were accounted worthy to suffer reproach for the name of Jesus.
  <b>42</b> And every day they ceased not in the temple, and from house to house, to teach and preach Christ Jesus.
</code></pre>
<pre>
<code>>>> scriptura Ephesians 6:10-
  <b>10</b> Finally, brethren, be strengthened in the Lord, and in the might of his power.
  <b>11</b> Put you on the armour of God, that you may be able to stand against the deceits of the devil.
  <b>12</b> For our wrestling is not against flesh and blood; but against principalities and power, against the rulers of the world of this darkness, against the spirits of wickedness in the high places.
... through verse 24 (end of chapter)
</code></pre>
<pre>
<code>>>> scriptura Revelation 22
  <b>1</b> And he showed me a river of water of life, clear as crystal, proceeding from the throne of God and of the Lamb.
  <b>2</b> In the midst of the street thereof, and on both sides of the river, was the tree of life, bearing twelve fruits, yielding its fruits every month, and the leaves of the tree were for the healing of the nations.
  <b>3</b> And there shall be no curse any more; but the throne of God and of the Lamb shall be in it, and his servants shall serve him.
... through verse 21 (end of chapter)
</code></pre>
<pre>
<code>>>> scriptura Mark 1-2
  <b>Chapter 1</b>
  <b>1</b> The beginning of the gospel of Jesus Christ, the Son of God.
  <b>2</b> As it is written in Isaias the prophet: Behold I send my angel before thy face, who shall prepare the way before thee.
... through chapter 2
</code></pre>
<pre>
<code>>>> scriptura Ecclesiastes
  <b>Chapter 1</b>
  <b>1</b> The words of Ecclesiastes, the son of David, king of Jerusalem.
  <b>2</b> Vanity of vanities, said Ecclesiastes vanity of vanities, and all is vanity.
... through chapter 12 (end of book)
</code></pre>

## Install

### Binary releases

Binaries for Windows and Linux on x86_64 are [available](https://github.com/john-bieren/scriptura/releases).

### Build from source

1. Clone the repository:
    ```
    git clone https://github.com/john-bieren/scriptura.git
    ```
2. Navigate to the project directory and compile the program:
    ```
    go build
    ```
