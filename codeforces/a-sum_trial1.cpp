#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;

int main() {
    int x;
    if (!(cin >> x) || x < 1 || x > 9261) {
        cout << "Invalid input" << endl;
        return 1;
    }
    for (int i = 0; i < x; ++i) {
        int a, b, c;
        if (!(cin >> a >> b >> c)) {
            cout << "Invalid input" << endl;
            return 1;
        }
        int max_val = max({a, b, c});
        int sum_val = a + b + c;
        if (2 * max_val == sum_val) {
            cout << "YES" << endl;
        } else {
            cout << "NO" << endl;
        }
    }

    return 0;
}