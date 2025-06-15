# Minimal Systemigram Description Language (MSDL) Specification Version 0.1

## 1. Overview

The Minimal Systemigram Description Language (MSDL) is a lightweight modeling language for expressing systems through structured systemigrams (SSGs). It prioritizes readability, simplicity, and pragmatic reasoning. MSDL is designed to:

- Be easy to learn and use

- Make diagrams self-explanatory

- Support minimal but powerful expressivity

- Provide immediate visual feedback, relying on layout algorithms for rendering

MSDL intentionally avoids elaborate features in favor of portability and clarity. 

Structured systemigrams are exactly the diagrams defined by this specification.

MSDL can also be used for systemigrams that are not structured, like originally devised by John Boardman and Brian Sauser in Systemic Thinking - Building Maps for Worlds of Systems, 2013, if Venn-type inclusion is replaced by structural edges.

**Structured** systemigrams specifcally adhere to the rules laid out in **3. Syntax Rules** and **4. Semantics and Interpretation**.

## 2. File Format

- MSDL files use the `.ssg` extension

- Files are plain text

- Files are expected be UTF-8 encoded.

## 3. Syntax Rules

### Title

- `{Title}` sets the title of the systemigram
- Only the first valid line matching `{...}` will be interpreted as the diagram title. Later occurences will be ignored.

### Nodes and Edges

- Behavioral edge: `A (verb) B.`

- Colored behavioral edge: `A (verb, color) B.`

- Structural edge: `A ((verb)) B.`

- Colored structural edge: `A ((verb, color)) B.`

- Node styling through colored border: `[Node]` (defaults to red) or `[Node, color]`

- Repeated identical edge declarations add edges, repeated identical node declarations do nothing. Repeated node declaration with different colors result undefined behavior.

- Special characters such as `{}`, `()`, `\`, and `"` are not supported inside node or verb names. 

### Comments

- `// Comment text` for single-line comments

- `/* ... */` for multi-line comments

- This facilitates literate descriptional diagramming

### Grammar Notes

- Statements should end with a period `.`

- Whitespace is ignored around tokens

- Node and verb phrases can include spaces

### Colors

- Color names must be Graphviz-compatible (e.g., `red`, `gray`, `blue`)

### Feature Completeness

- MSDL includes only the syntax features listed above. No additional syntax is supported.

## 4. Semantics and Interpretation

### Behavioral Edges

Express what *happens* - they capture system behavior. The semantics are pragmatically specified as follows:

- `A (verb describing change) B`. This is understood as:  A influences B’s change. A itself remains in this interaction relatively unchanged, while B undergoes the interesting transformation. The verb can be in active or passive voice depending on the direction of the change and the verb used. Compare `sensor (activates) alarm` to `environment (is scanned by) scanner`. In the later case the scanner's state is changed, while the environment is unchanged.  

- `A (outputs an object to) B.` `A (is input into) B.` This is understood as:  A inputs something, that is then changed in B, probably output again. The changed object does not have to be an explicit node anywhere. It can be only tracked through stations that change it. Or the object is modeled explicitly, A is an input to B, that is then used or changed in B. A is then often produced by some upstream node.

- Behavioral edges point to the object which is changed (or changed the most) through the behavior, or where an implicit or explicit other object is transported to. The sentence form follows the edge: subject --(predicate)--> object.

- Behavioral edges stand for classes of behavior. **Instances** of behavioral edges take time, and run one after the other. They can run in the background in multiple instances while you inspect one of them, or even continuously.

### Structural Edges

Describe *stable* relationships, that are not behavior.

- `A ((verb describing connection)) B.` This is understood as:  A is connected to B, and this connection does not change.  Examples are `A is part of B.`, `A owns B.`, `A hosts B.`
- Structural edges don't run and do not take time. They are there.

### Colors and Extensions

- Modelers may use colors to create semantic extensions. They must be explained in a comment in the MSDL file, to ensure that the model is understandable on its own. Colors without explanation amount to highlighting (e.g. for the purpose of a presentation).

## 5. Examples

```msl
{Simple Heating System}
oven (produces, red) heat.
heat (raises) temperature.
temperature (is checked by, blue) sensor.
sensor (regulates) oven.
sensor ((is part of)) wall.
wall ((contains, gray)) oven.
[oven]
[sensor, green]
```

The colors do not have any set meaning in this example.

Multi-line Markdown annotation:

```msl
/**********
# Simple Heating Example
- Red arrow: oven produces heat
- Blue arrow: sensor checks temperature
- Green node: sensor is highlighted
**********/
```

This can go before, in the middle, or after other declarations. 

## 6. Tooling Assumptions

- Editors parse the file line-by-line

- First `{Title}` is taken as the title and shown atop the rendered diagram

- Multi-line and single-line comments are stripped before parsing statements

- Node and edge declarations are translated into Graphviz DOT format

- The dot diagram is displayed to the user in real time

## 7. Translation to Graphviz

MSDL statements are translated into Graphviz DOT using the following rules:

- All systemigrams are translated into a Graphviz `digraph` block:
  
  ```
  digraph G {
  ...
  }
  ```
  
  This means that all edges are treated as **directed** by default, reflecting the behavioral or structural directionality of systemigram relations.

- Each node in an edge statement becomes a DOT node, e.g. `oven (produces) heat.` => `oven -> heat [label="produces"];`

- Colored edges compile as: 
  `A (verb, red) B.` => `A -> B [label="verb", color="red"];`

- Structural edges are rendered with `style=dashed` and optional color: 
  `A ((hosts, gray)) B.` => `A -> B [label="hosts", style=dashed, color="gray"];`

- Node styling uses `color` attributes: 
  `[oven]` => `oven [color="red"];`  (default)
  `[sensor, green]` => `sensor [color="green"];`

- The first `{Title}` line is mapped to a top-level `label` directive: 
  `{My System}` => `label="My System"; labelloc=top; fontsize=20;`

- Comment lines (`//`) and comment blocks (`/* ... */`) are ignored by the compiler.

## 8. Non-goals

MSDL deliberately avoids any complexity beyond what is needed to create structured systemigrams to reason about and communicate system behavior. *"In the end, remedying the SoI is really all that matters."* (John Boardman and Brian Sauser)

MSDL is not an MBSE language, and cannot be, thought it took inspiration from UML, SysML and OPM. Its role is to explore and explain systems.

For advanced styling or automation, export to and post-process in Graphviz.

## 9. Feature Policy

MSDL exists to lower the entry barrier to systems thinking and modeling. It achieves this by being simple and small.

Feature requests must demonstrate that the proposed addition would measurably improve reasoning or communication across a broad range of use cases for non-sophisticated users. The bar is intentionally high to preserve MSDL’s primary goal.

<!--
Copyright 2025 Stefan Kober

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

A copy of the License is also included in this repository.

In short: Preserve credit. If you change this file, mark your changes clearly.
-->
