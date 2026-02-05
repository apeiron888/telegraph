import math
for _ in range(int(input())):
    n = int(input())
    lst = list(map(int,input().split()))
    gcf = lst[0]
    for i in range(1,n):
        gcf = math.gcd(gcf,lst[i])
    for j in range(2,gcf):
        if math.gcd(j,gcf)==1:
            print(j)
            break
    else:
        print(gcf+1 if gcf<10**18 else -1)