class Solution:
    def nextGreaterElement(self, nums1: List[int], nums2: List[int]) -> List[int]:
        counter = {nums2[idx]:idx for idx in range(len(nums2))}
        stack = [nums2[0]]
        for idx in range(1,len(nums2)):
            while stack and nums2[idx] > stack[-1]:
                nums2[counter[stack[-1]]]=nums2[idx]
                stack.pop()
            stack.append(nums2[idx])
        while stack:
            nums2[counter[stack[-1]]] = -1
            stack.pop()
        for idx in range(len(nums1)):
            nums1[idx] = nums2[counter[nums1[idx]]]
        return nums1
         