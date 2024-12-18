# Spec

## Body

It's based on a subset of CommonMark Markdown

This is very similar to the format used by [Grocery](https://github.com/cnstoll/Grocery-Recipe-Format), likely close enough that Grocery files will work as RecipeMark, and mostly vice-versa

```
# Title
## Header
### Ingredient or recipe subheader
- Ingredient
Step
> Note
```

An ingredient can have up to three sections, separated by pipe characters. Ingredient, quantity, preparation or notes
```
- Sugar | 1tsp | for dusting
```

Tags can be used inline as #tag, and links to other recipes handled via [[wikilinks]]

## Metadata

Metadata is added in YAML format in the head of the document in the usual Markdown way.

Supported metadata, all optional, includes

 * name - defaults to first # Header
 * description - defaults to first blockquote
 * image - the name of the header image
 * datePublished - defaults to file timestamp
 * author
 * prepTime
 * cookTime
 * totalTime
 * cookingMethod
 * category
 * cuisine
 * yield