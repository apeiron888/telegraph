def findUnion(a, b):
    # Convert both arrays to sets to remove duplicates
    # Use the union operator '|' to combine them
    union_set = set(a) | set(b)
    
    # Return as a list (the driver code will handle sorting)
    return list(union_set)