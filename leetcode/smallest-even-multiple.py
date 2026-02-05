class Solution:
    def smallestEvenMultiple(self, n: int) -> int:
        # LCM * GCF = n*2 
        # GCF = 1 if n is odd and GCF = 2 if n is even 
        return n if n%2==0 else 2*n