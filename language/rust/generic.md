# Generic types
1. Generic structure
2. Generic function
3. Generic trait

# trait
Define a group of structure which share common behavior.
Normally used in two ways:
1. As constraint as typeclass in haskell to make sure generic
   type parameter belongs to some group/kind
2. As trait object which provide dynamic dispatch

# Associate type vs generic trait
1. Generic trait is a way to define specific trait, different
   type parameter defines different trait 
2. Associate type in trait is another implementation/behavior 
   detail specified for the structure which belongs to this trait.

# Basic idea
Use other types as parameter to generate concrete type
Could set constraint for the type parameter use trait
