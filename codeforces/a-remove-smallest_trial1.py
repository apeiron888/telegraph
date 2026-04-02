import sys
from typing import List

def main():
    t = int(input())
    for _ in range(t):
        n = int(input())
        a = list(map(int, input().split()))
        if n == 1:
            print("YES")
            continue
        a.sort()
        v = []
        for i in range(1, n):
            dif = abs(a[i] - a[i-1])
            v.append(dif)
        v.sort(reverse=True)
        if v[0] > 1:
            print("NO")
        else:
            print("YES")

if __name__ == "__main__":
    main()
