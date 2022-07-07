# Client API controller
Not usually a fan of the whole redux kitchen sink, it would be nice to cut down some request boilerplate for retrieving stuff.

## benefits
- enable \[result, error] pattern
- organization
- die and do nothing when the navigation needs to be hijacked

## example
```javascript
// controller thing
// success: [[{Name: "foo"}], false]
// actual error: [null, new Error("something")]
// die: [null, true]
export const listCharacters = () => {...};

// consumer
const [characters, error] = await listCharacters();
```