class Solution:
    def validPalindrome(self, s: str) -> bool:
        left, right, chance = 0, len(s)-1, 1

        while left < right:
            if s[left]!=s[right]:
                possibility_1 = s[left+1:right+1]
                possibility_2 = s[left:right]
                return possibility_1==possibility_1[::-1] or possibility_2 == possibility_2[::-1]        
            left += 1
            right -= 1
    
        return True