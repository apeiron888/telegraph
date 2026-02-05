#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

bool validate(int mid, int m, const vector<int>& street, const vector<int>& g_count) {
    for (int i = 0; i < 26; ++i) {
        if (street[i] == 0) continue;
        long long rest = (long long)g_count[i] - street[i];
        if (rest == 0) return false;
        if (1LL * (m - mid) * rest < (long long)street[i]) {
            return false;
        }
    }
    return true;
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);

    int n, m;
    if (!(cin >> n >> m)) return 0;

    vector<string> lst(n);
    vector<int> g_count(26, 0);
    vector<vector<int>> c_counts(n, vector<int>(26, 0));

    for (int i = 0; i < n; ++i) {
        cin >> lst[i];
        for (char c : lst[i]) {
            int idx = c - 'A';
            g_count[idx]++;
            c_counts[i][idx]++;
        }
    }
    for (int i = 0; i < n; ++i) {
        int low = 0, high = m - 1;
        int ans_k = -1;

        while (low <= high) {
            int mid = low + (high - low) / 2;
            if (validate(mid, m, c_counts[i], g_count)) {
                ans_k = mid;
                low = mid + 1;
            } else {
                high = mid - 1;
            }
        }

        cout << ans_k << (i == n - 1 ? "" : " ");
    }
    cout << endl;

    return 0;
}