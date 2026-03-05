cr, cc = 3, 3
for i in range(1,6):
    row = list(map(int, input().split()))
    if 1 in row:
        j = row.index(1) + 1
        move = abs(i-cr)+abs(j-cc)
        print(move)
        