for _ in range(int(input())):
    n = int(input())
    a = list(map(int,input().split()))
    total = sum(a)
    ones = a.count(1)
    if n > 1 and total >= n + ones:
        print("YES")
    else:
        print("NO")
        