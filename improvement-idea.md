# Dokumen Ide Peningkatan (Improvement Backlog)

Dokumen ini mencatat ide-ide penyempurnaan, variasi mekanik, dan peningkatan kualitas permainan (quality of life) yang dapat diimplementasikan di kemudian hari setelah fondasi utama tiap milestone stabil.

---

## 1. Fluctuating Commodity Yields & Dynamic Weather System

* **Bagian / Sub-Sistem:** `Milestone 1: Settlement Simulation & Worker Engine` ([internal/settlement/simulation.go](file:///media/fahmi/data/utility/terminal-odyssey/internal/settlement/simulation.go))
* **Status Saat Ini:** Masih menggunakan model deterministik murni (produksi harian bernilai tetap sesuai jumlah pekerja dan pengali musim).
* **Tujuan Peningkatan:** Memberikan dinamika kejutan harian yang lebih hidup dan realistis tanpa menghilangkan kemampuan pemain dalam merencanakan alokasi pekerja secara strategis.

### Rincian Konsep Peningkatan:

#### A. Variansi Hasil Harian (Daily Yield Variance & Critical Harvest)
* **Deviasi Acak Wajar:** Menambahkan rentang fluktuasi kecil ($\pm 10\% - 20\%$) pada hasil kerja harian tiap peran pekerja.
  * *Contoh:* 1 Petani menghasilkan `3 – 5` ransum (rata-rata 4), alih-alih selalu pasti 4 setiap hari.
* **Mekanik *Critical Harvest*:** Memberikan peluang kecil (~5% - 10%) bagi pekerja untuk memicu kejadian luar biasa:
  * *Panen Raya:* Ladang menghasilkan bonus $+50\%$ ransum hari itu.
  * *Urat Bijih Murni:* Penambang menemukan kantong mineral padat yang memberikan bonus batu atau bijih besi ekstra.

#### B. Sistem Cuaca Harian (Dynamic Daily Weather)
Sistem cuaca harian yang berjalan di bawah payung 4 musim besar (Musim menentukan probabilitas cuaca yang muncul):
* **Hujan Lebat / Badai:**
  * Meningkatkan kesuburan tanah (Panen $+25\%$).
  * Menggenangi lubang galian tambang (Produksi Batu $-30\%$).
* **Kemarau Terik:**
  * Mengeringkan sumber air (Produksi Ransum $-20\%$).
  * Mempermudah transportasi dan pengeringan kayu (Penebangan $+15\%$).
* **Kabut Tebal:**
  * Jarak pandang terbatas, efisiensi penebang hutan berkurang $-15\%$.
  * Garda Milisi membutuhkan kewaspadaan ganda terhadap penyusupan bandit.
* **Badai Salju Ekstrem (Khusus Musim Dingin):**
  * Melipatgandakan kebutuhan kayu bakar per hari (+50% konsumsi kayu).
