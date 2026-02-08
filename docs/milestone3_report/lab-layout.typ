#let title-page(
  title:[Title Here], 
  sub: [], 
  email:[], 
  osid:[], 
  author: "Edgar J Coime", 
  description: [],
  running-head: none,
  font-family: "Libertinus Serif", 
  font-size: 12pt,
  region: "us",
  language: "en",
  paper-size: "us-letter",
  fill: white, 
  body
) = {
  let double-spacing = 1.5em
  let first-indent-length = 0.5in

  set document(
    title: title,
    author: author,
    description: description,
  )

  // TITLE PAGE
  set page(fill: rgb(fill), margin: (top: 1.5in, rest: 2in), paper: paper-size)
  set text(font: font-family, size: 12pt, weight: "regular")
  align(left)[
    British Columbia Institute of Technology
    #v(0.1em)
    #text(size: 10pt, "Network Security Applications Development")
  ]
  line(start: (0%, 0%), end: (8.5in, 0%), stroke: (thickness: 4pt))
  align(horizon + left)[
    #text(size: 24pt, title, weight: "medium")\
    #v(1em)
    #text(sub)
    #v(2em)
    #author \
    #osid \
    #link(email)
  ]
  
  align(bottom + left)[#datetime.today().display()]
  pagebreak()

  // Document settings
  set text(
    size: font-size,
    font: font-family,
    region: region,
    lang: language,
  )

  set page(margin: 1in, paper: paper-size, numbering: "1", number-align: top + right, header: context {
    upper(running-head)
    h(1fr)
    str(here().page())
  })

  set par(
    leading: double-spacing,
    spacing: double-spacing,
  )

  show link: it => {
    if type(it.dest) == str {
      set text(fill: blue)
      underline(it)
    } else { it }
  }

  if running-head != none {
    if type(running-head) == content { running-head = to-string(running-head) }
    if running-head.len() > 50 {
      panic("Running head must be no more than 50 characters, including spaces and punctuation.")
    }
  }

  show heading: set text(size: 14pt)
  show heading: set block(spacing: double-spacing)

  show heading: it => emph(strong[#it.body.])
  show heading.where(level: 1): it => align(center, strong(it.body))
  show heading.where(level: 2): it => par(first-line-indent: 0in, strong(it.body))

  show heading.where(level: 3): it => par(first-line-indent: 0in, emph(strong(it.body)))

  show heading.where(level: 4): it => strong[#it.body.]
  show heading.where(level: 5): it => emph(strong[#it.body.])

  set par(
    first-line-indent: (
      amount: first-indent-length,
      all: true,
    ),
    leading: double-spacing,
  )

//  show table.cell: set par(leading: 1em)

//  show figure: set block(breakable: true, sticky: false)
//
//  set figure(
//    gap: double-spacing,
//    placement: auto,
//  )
//
//  set figure.caption(separator: parbreak(), position: top)
//  show figure.caption: set align(left)
//  show figure.caption: set par(first-line-indent: 0em)
//  show figure.caption: it => {
//    strong[#it.supplement #context it.counter.display(it.numbering)]
//    it.separator
//    emph(it.body)
//  }

//  set table(stroke: (x, y) => if y == 0 {
//    (
//      top: (thickness: 1pt, dash: "solid"),
//      bottom: (thickness: 0.5pt, dash: "solid"),
//    )
//  })

  set list(
    marker: ([•], [◦]),
    indent: 0.5in - 1.75em,
    body-indent: 1.3em,
  )

  set enum(
    indent: 0.5in - 1.5em,
    body-indent: 0.75em,
  )

  set raw(
    tab-size: 4,
    block: true,
  )

  show raw.where(block: true): block.with(
    fill: luma(250),
    stroke: (left: 3pt + rgb("#6272a4")),
    inset: (x: 10pt, y: 8pt),
    width: auto,
    breakable: true,
    outset: (y: 7pt),
    radius: (left: 0pt, right: 6pt),
  )

  show raw: set text(
    font: "Cascadia Code",
    size: 10pt,
  )

  show raw.where(block: true): set par(leading: 1em)
  show figure.where(kind: raw): set block(breakable: true, sticky: false, width: 100%)

  set math.equation(numbering: "(1)")

  show quote.where(block: true): set block(spacing: double-spacing)

  show quote: it => {
    let quote-text-words = to-string(it.body).split(regex("\\s+")).filter(word => word != "").len()

    if quote-text-words < 40 {
      ["#it.body" ]

      if (type(it.attribution) == label) {
        cite(it.attribution)
      } else if (
        type(it.attribution) == str or type(it.attribution) == content
      ) {
        it.attribution
      }
    } else {
      block(inset: (left: 0.5in))[
        #set par(first-line-indent: 0.5in)
        #it.body
        #if (type(it.attribution) == label) {
          cite(it.attribution)
        } else if (type(it.attribution) == str or type(it.attribution) == content) {
          it.attribution
        }
      ]
    }
  }

  show outline.entry: it => {
    if (
      (
        it.element.supplement == [#context get-terms(language).Appendix]
          or it.element.supplement == [#context get-terms(language).Annex]
          or it.element.supplement == [#context get-terms(language).Addendum]
      )
        and it.element.has("level")
        and it.element.level == 1
    ) {
      link(it.element.location(), it.indented([#it.element.supplement #it.prefix().], it.inner()))
    } else {
      it
    }
  }

  set outline(depth: 3, indent: 2em)

  set bibliography(style: "apa")
  show bibliography: set par(
    first-line-indent: 0in,
    hanging-indent: 0.5in,
  )

  show bibliography: bib-it => {
    set block(inset: 0in)
    show block: block-it => context {
      if block-it.body.func() != [].func() {
        block-it.body
      } else {
        par(block-it.body)
      }
    }

    bib-it
  }
  body
}
