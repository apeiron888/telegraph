try:
    x = int(input())
    if x<1 or x>9261:
        print("Invalid input")
        exit(1)
    for _ in range(x):
        a,b,c = [int(x) for x in input().split()]
        print( "YES" if 2*max(a,b,c) == a+b+c else "NO")
except:
    print("Invalid input")
    exit(1)