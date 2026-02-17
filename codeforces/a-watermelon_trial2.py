try:
    x = int(input())
    if x<1 or x>100:
        print("Invalid number: 1<=x<=100")
        exit(1)
    if x%2==0 and x!=2: print("YES")
    else: print("NO")
    
except ValueError:
    print("Invalid input: only numbers allowed")
