# Terminal Odyssey

[![CI](https://github.com/Fahmi-mi/terminal-odyssey/actions/workflows/ci.yml/badge.svg)](https://github.com/Fahmi-mi/terminal-odyssey/actions/workflows/ci.yml)
[![Release](https://github.com/Fahmi-mi/terminal-odyssey/actions/workflows/release.yml/badge.svg)](https://github.com/Fahmi-mi/terminal-odyssey/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Fahmi-mi/terminal-odyssey)](https://golang.org)

Terminal Odyssey adalah game Command-Line Interface (CLI) modern berbasis teks yang memadukan simulasi manajemen pemukiman (*settlement management*), sistem ekonomi dan perdagangan komoditas (*trading simulation*), serta penjelajahan bawah tanah naratif prosedural (*narrative roguelike dungeon crawler*).

Game ini dibangun menggunakan bahasa pemrograman **Go (Golang)** dengan memanfaatkan framework Text User Interface (TUI) **Bubble Tea** dan sistem tata letak serta pewarnaan **Lip Gloss**.

---

## Daftar Isi

- [Latar Belakang & Premis](#latar-belakang--premis)
- [Fitur Utama](#fitur-utama)
- [Persyaratan Sistem](#persyaratan-sistem)
- [Cara Instalasi & Menjalankan Permainan](#cara-instalasi--menjalankan-permainan)
- [Perintah Makefile](#perintah-makefile)
- [Panduan Navigasi & Pintasan Keyboard](#panduan-navigasi--pintasan-keyboard)
- [Panduan Mekanika Permainan](#panduan-mekanika-permainan)
  - [1. Manajemen Pemukiman & Siklus Harian](#1-manajemen-pemukiman--siklus-harian)
  - [2. Empat Musim Dinamis](#2-empat-musim-dinamis)
  - [3. Penjelajahan Katakombe & Pertempuran Taktis](#3-penjelajahan-katakombe--pertempuran-taktis)
  - [4. Bengkel Pandai Besi & Pelatihan Stat](#4-bengkel-pandai-besi--pelatihan-stat)
  - [5. Pasar Komoditas & Serikat Kafilah](#5-pasar-komoditas--serikat-kafilah)
  - [6. Laboratorium Alkimia & Kedai Minum](#6-laboratorium-alkimia--kedai-minum)
  - [7. Pengepungan Gerbang Pemukiman](#7-pengepungan-gerbang-pemukiman)
  - [8. Boss Katakombe & Kondisi Kemenangan](#8-boss-katakombe--kondisi-kemenangan)
  - [9. Sistem Simpan & Muat Permainan](#9-sistem-simpan--muat-permainan)
- [Struktur & Arsitektur Proyek](#struktur--arsitektur-proyek)
- [Otomasi CI/CD & Rilis](#otomasi-cicd--rilis)
- [Pengujian Kode](#pengujian-kode)
- [Lisensi](#lisensi)

---

## Latar Belakang & Premis

Sebagai pemimpin pemukiman terpencil bernama **Oakhaven**, Anda memikul tanggung jawab ganda:
1. Membangun perkemahan darurat menjadi kota benteng dagang yang makmur melalui pengelolaan warga, alokasi logistik, penempaan senjata, peracikan obat, dan pengiriman kafilah niaga.
2. Memimpin ekspedisi ke kedalaman **Katakombe Abyssal** kuno untuk merebut kembali relik purba, menghadapi monster bayangan, menguji kewarasan mental di tengah kegelapan, dan menumbangkan ancaman **Abyssal Overlord** demi mengamankan masa depan Oakhaven.

---

## Fitur Utama

- **Simulasi Pemukiman & Pekerja Real-Time Turn:** Kelola alokasi Petani, Penebang Kayu, Penambang Batu, Pandai Besi, dan Garda Milisi. Bangun dan tingkatkan 9 jenis fasilitas desa.
- **Katakombe Prosedural Berbasis Pacing:** Struktur bawah tanah dinamis dengan sistem obor (*Torch*), tingkat stres kewarasan (*Sanity*), trauma psikologis, pertemuan acak, ruang misteri, dan suaka pemulihan.
- **Pertempuran Taktis Berbasis Giliran (Turn-Based Combat):** Hadapi musuh katakombe dengan mekanik serangan, bertahan, penggunaan perbekalan ransum, dan minum ramuan alkimia tempur.
- **Tempa Senjata Modular & Progresi Karakter:** 4 kelas senjata (Pedang, Belati, Gada, Tombak) dengan efek affix khusus dan durabilitas. Latih 4 stat utama: *Might*, *Agility*, *Resolve*, dan *Ingenuity*.
- **Pasar Bebas & Ekspedisi Kafilah Dagang:** Fluktuasi harga beli/jual 5 komoditas (Gandum, Kulit Hewan, Besi Tempa, Sutra Gurun, Rempah Kuno) serta ekspedisi kafilah regional berisiko tinggi dengan margin arbitrase besar.
- **Laboratorium Alkimia & Perekrutan Pendamping:** Racik ramuan di laboratorium apotek desa dan rekrut 4 spesialisasi tentara bayaran di kedai minum (*Vanguard*, *Rogue*, *Scholar*, *Acolyte*) yang memberikan kemampuan pasif saat bertualang.
- **Pengepungan Pemukiman (Siege Engine):** Akumulasi kekayaan memicu ancaman penjarah. Garda milisi, menara pengawas, dan benteng pertahanan bertarung otomatis di gerbang desa saat pergantian hari.
- **Pertarungan Boss Akhir & Layar Kemenangan Naratif:** Tantang Abyssal Overlord di kedalaman lantai 5. Penuhi 3 syarat kemenangan untuk menamatkan skenario naratif atau lanjutkan dalam mode sandbox tanpa batas.
- **Persistensi Simpan & Muat (Save/Load):** Format data JSON mandiri dengan 3 slot penyimpanan manual serta 1 slot penyimpanan otomatis (*autosave*).
- **Antarmuka CLI Terstandarisasi:** Tampilan 78 kolom konsisten di setiap layar dengan palet warna TrueColor yang nyaman di mata.

---

## Persyaratan Sistem

- **Sistem Operasi:** Linux, macOS, atau Windows.
- **Bahasa Pemrograman:** Go versi 1.21 ke atas (pengembangan menggunakan Go 1.26.5).
- **Terminal Emulator:** Terminal modern dengan dukungan UTF-8, ANSI escape sequences, dan TrueColor (minimal ukuran jendela terminal 80 kolom x 24 baris).
- **Build Tool:** GNU Make (opsional, untuk menjalankan target otomatis).

---

## Cara Instalasi & Menjalankan Permainan

### 1. Kloning Repositori

```bash
git clone https://github.com/Fahmi-mi/terminal-odyssey.git
cd terminal-odyssey
```

### 2. Menjalankan Langsung via Go CLI

```bash
go run ./cmd/odyssey
```

### 3. Kompilasi Binary Lokal

```bash
# Menggunakan Go CLI langsung
go build -ldflags="-s -w" -o bin/odyssey ./cmd/odyssey

# Menjalankan binary yang dihasilkan
./bin/odyssey
```

---

## Perintah Makefile

Tersedia `Makefile` untuk mempermudah eksekusi alur kerja harian:

| Perintah | Deskripsi |
|---|---|
| `make run` | Menjalankan game secara langsung tanpa kompilasi manual |
| `make build` | Mengompilasi binary tunggal ke `bin/odyssey` |
| `make test` | Menjalankan seluruh rangkaian unit test di semua paket |
| `make test-coverage` | Menjalankan unit test dan menampilkan persentase cakupan kode |
| `make vet` (atau `make lint`) | Memeriksa kepatuhan kode statis menggunakan `go vet` |
| `make clean` | Menghapus direktori `bin/` dan berkas profil pengujian |
| `make build-all` | Melakukan *cross-compilation* binary untuk Linux, Windows, dan macOS |
| `make all` | Menjalankan `vet`, `test`, dan `build` secara berurutan |
| `make help` | Menampilkan ringkasan bantuan perintah Makefile |

---

## Panduan Navigasi & Pintasan Keyboard

Game ini tidak menggunakan pergerakan ubin grid/WASD yang kaku, melainkan menggunakan navigasi berbasis pilihan menu dan tombol pintasan (*shortcuts*).

### Menu Utama (Title Screen)
- `[1]` / `[N]`: Memulai permainan baru (*New Game*)
- `[2]` / `[L]`: Membuka menu muat permainan (*Load Game*)
- `[3]` / `[Q]`: Keluar dari permainan (*Quit*)
- `[Panah Atas / Bawah]` atau `[k / j]`: Memilih menu atau memilih slot simpanan
- `[Enter]`: Konfirmasi pilihan menu atau muat slot tersimpan
- `[Esc]` / `[B]`: Kembali ke layar awal dari daftar slot

### Markas Pemukiman (Town Menu)
- `[1]`: Masuk ke Pasar & Perdagangan Komoditas
- `[2]`: Masuk ke Pembangunan & Peningkatan Fasilitas Desa
- `[3]`: Masuk ke Bengkel Pandai Besi (Tempa Senjata & Armory) *(Butuh Bengkel)*
- `[4]`: Masuk ke Pusat Latihan Stat Karakter *(Butuh Pusat Latihan)*
- `[5]`: Masuk ke Kedai Minum (Rekrut Pendamping & Rumor) *(Butuh Kedai)*
- `[6]`: Berangkat Menuju Katakombe Bawah Tanah (Membawa Ransum)
- `[7]`: Masuk ke Laboratorium Alkimia (Racik Ramuan) *(Butuh Apotek)*
- `[W]`: Buka Pengaturan Alokasi Tenaga Kerja Warga
- `[D]`: Lewati Hari (Jalankan siklus produksi & konsumsi harian)
- `[S]`: Buka Menu Simpan Permainan (*Save Game*)
- `[M]`: Kembali ke Layar Judul (*Title Screen*)
- `[Q]` / `[Ctrl+C]`: Keluar dari permainan

### Alokasi Tenaga Kerja (Worker Assign)
- `[Panah Atas / Bawah]`: Memilih peran pekerja (Petani, Penebang, Penambang, Pandai Besi, Milisi)
- `[Panah Kanan]` / `[+]`: Menugaskan +1 warga ke peran yang dipilih
- `[Panah Kiri]` / `[-]`: Menarik 1 pekerja dari peran yang dipilih
- `[Esc]`: Selesai dan kembali ke Markas Pemukiman

### Pembangunan Pemukiman (Building Construction)
- `[Panah Atas / Bawah]` atau Angka `[1]` - `[9]`: Memilih fasilitas pemukiman
- `[Enter]`: Meningkatkan level fasilitas yang dipilih (menggunakan material gudang)
- `[Esc]`: Kembali ke Markas Pemukiman

### Bengkel Pandai Besi (Blacksmith)
- `[Tab]` / Angka `[1]` - `[2]`: Berpindah antara Tab Tempa Baru dan Tab Armory Senjata
- `[Panah Atas / Bawah]`: Memilih resep tempa atau memilih senjata koleksi
- `[Enter]`: Tempa senjata baru atau gunakan (*equip*) senjata yang dipilih
- `[R]`: Memperbaiki durabilitas senjata yang aus menggunakan bijih batu
- `[Esc]`: Kembali ke Markas Pemukiman

### Pusat Latihan Karakter (Training Grounds)
- `[Panah Atas / Bawah]`: Memilih stat inti (*Might*, *Agility*, *Resolve*, *Ingenuity*)
- `[Enter]`: Meningkatkan stat yang dipilih (menggunakan emas dan ransum desa)
- `[Esc]`: Kembali ke Markas Pemukiman

### Pasar & Kafilah Dagang (Market & Caravans)
- `[Tab]` / Angka `[1]` - `[2]`: Berpindah antara Tab Pasar Lokal dan Tab Ekspedisi Kafilah
- `[Panah Atas / Bawah]`: Memilih komoditas atau memilih rute ekspedisi
- `[B]`: Membeli 1 unit komoditas yang dipilih
- `[S]`: Menjual 1 unit komoditas yang dipilih
- `[Enter]` (pada Tab Kafilah): Mengirim kafilah dagang ke rute yang dipilih
- `[Esc]`: Kembali ke Markas Pemukiman

### Kedai Minum (Tavern)
- `[Tab]` / Angka `[1]` - `[2]`: Berpindah antara Tab Rekrutmen Pendamping dan Tab Istirahat/Rumor
- `[Panah Atas / Bawah]`: Memilih tentara bayaran di dalam daftar kedai
- `[Enter]`: Rekrut atau berhentikan pendamping yang dipilih
- `[R]` (pada Tab Rumor): Beristirahat di kedai minum untuk memulihkan Sanity
- `[M]` (pada Tab Rumor): Mendengarkan desas-desus pedagang musafir
- `[Esc]`: Kembali ke Markas Pemukiman

### Laboratorium Alkimia (Alchemy Lab)
- `[Panah Atas / Bawah]`: Memilih resep ramuan
- `[Enter]`: Meracik ramuan yang dipilih ke dalam kantong petualang
- `[Esc]`: Kembali ke Markas Pemukiman

### Eksplorasi Katakombe Bawah Tanah (Dungeon Exploration)
- `[Enter]`: Membuka peti harta karun, beristirahat di suaka, atau menghadapi altar misteri
- `[1]`, `[2]`: Memilih cabang rute lorong menuju ruangan berikutnya
- `[M]`: Mengonsumsi 1 ransum ransel untuk memulihkan HP
- `[O]`: Menyalakan obor cadangan (+30% daya nyala obor)
- `[P]`: Menggunakan ramuan darurat yang ada di kantong petualang
- `[Esc]`: Melakukan evakuasi mundur kembali ke desa secara aman

### Pertempuran Taktis (Turn-Based Combat)
- `[1]`: Menyerang musuh dengan senjata aktif (*Attack*)
- `[2]`: Mengambil posisi bertahan (*Guard*, mereduksi 50% kerusakan musuh)
- `[3]`: Memulihkan HP darurat menggunakan ransum bekal (+25 HP)
- `[4]`: Mencoba melarikan diri dari medan laga (*Flee*)
- `[P]`: Meminum ramuan tempur (Salep Pemulih / Eliksir Kekuatan / Tonik / Penawar)
- `[Enter]`: Mengonfirmasi akhir duel dan melanjutkan penjelajahan

### Layar Laporan Pengepungan (Siege Report)
- `[Enter]` / `[Space]` / `[Esc]`: Menutup laporan pengepungan dan kembali ke desa

### Menu Simpan Permainan (Save Menu)
- `[Panah Atas / Bawah]` atau `[k / j]`: Memilih slot simpanan (`slot_1`, `slot_2`, `slot_3`)
- `[Enter]`: Menyimpan seluruh progres permainan ke slot yang dipilih
- `[Esc]`: Batal dan kembali ke desa

### Layar Kemenangan (Victory Screen)
- `[Enter]` / `[Space]`: Lanjutkan bermain dalam Mode Sandbox Tanpa Batas (*Endless Mode*)
- `[S]`: Simpan permainan status tamat
- `[Q]` / `[Esc]`: Kembali ke Menu Utama

---

## Panduan Mekanika Permainan

### 1. Manajemen Pemukiman & Siklus Harian
- **Populasi & Kapasitas:** Populasi warga dibatasi oleh level Balai Desa. Warga yang menganggur tidak menghasilkan sumber daya.
- **Konsumsi Ransum:** Setiap warga dan petualang mengonsumsi 1 ransum setiap pergantian hari. Pastikan jumlah Petani memadai agar lumbung pangan tidak mengalami defisit.
- **Kapasitas Gudang:** Kayu dan batu memiliki batas tampung yang ditentukan oleh level Gudang Logistik (*Storehouse*).

### 2. Empat Musim Dinamis
Siklus berjalan berurutan setiap 30 hari:
- **Musim Semi (Spring):** Efisiensi bibit naik, peluang migrasi warga baru meningkat.
- **Musim Panas (Summer):** Hasil panen melimpah (+35%), namun kebutuhan konsumsi air/ransum di dungeon meningkat.
- **Musim Gugur (Autumn):** Produksi penebangan kayu melonjak (+25%), persiapan bahan bakar menjelang musim dingin.
- **Musim Dingin (Winter):** Ladang pertanian membeku (produksi ransum terhenti total); desa memerlukan konsumsi kayu bakar harian.

### 3. Penjelajahan Katakombe & Pertempuran Taktis
- **Daya Nyala Obor (Torch):** Berkurang seiring langkah penjelajahan. Jika obor redup atau padam, akurasi serangan menurun dan risiko sergapan monster melonjak drastis.
- **Kewarasan Mental (Sanity):** Berkurang saat menghadapi teror kegelapan. Jika Sanity turun di bawah 30, karakter berisiko menderita trauma psikologis (*Paranoia*, *Klaustrofobia*, atau *Megalomania*).
- **Pendamping Petualang:** Pendamping yang direkrut memberikan efek pasif aktif di dungeon (contoh: *Rogue* melucuti jebakan, *Vanguard* menyerap kerusakan tim, *Scholar* menggandakan temuan relik, *Acolyte* memulihkan stres).

### 4. Bengkel Pandai Besi & Pelatihan Stat
- **Tipe Senjata:**
  - *Pedang Baja:* Gaya seimbang dengan peluang kritikal stabil.
  - *Belati Cepat:* Menyerang cepat dengan efek racun (*Poison*) berkala.
  - *Gada Berat:* Serangan tumpul berdaya rusak tinggi dengan peluang membuat musuh pingsan (*Stun*).
  - *Tombak Panjang:* Memberikan inisiatif serangan pendahuluan sebelum musuh mendekat.
- **Atribut Karakter:**
  - *Might:* Meningkatkan daya rusak fisik serangan senjata melee.
  - *Agility:* Meningkatkan peluang menghindar (*Dodge*) dan serangan kritikal.
  - *Resolve:* Memperbesar batas HP maksimal dan memperlambat laju kegilaan Sanity.
  - *Ingenuity:* Membuka efisiensi ramuan alkimia dan peluang tawar-menawar harga pasar.

### 5. Pasar Komoditas & Serikat Kafilah
- **Pasar Lokal:** Harga 5 komoditas bergerak naik/turun setiap hari berdasarkan penawaran, permintaan, dan musim regional.
- **Ekspedisi Kafilah:** Mengirim gerobak barang ke wilayah tetangga (*Pelabuhan Garam*, *Kota Benteng Besi*, *Oasis Pasir Hitam*). Risiko serangan penyamun berbanding lurus dengan potensi keuntungan emas yang dibawa pulang.

### 6. Laboratorium Alkimia & Kedai Minum
- **Ramuan Alkimia:** Mengolah bahan herba dan mineral dungeon menjadi *Salep Pemulih* (regenerasi HP), *Minyak Obor* (daya tahan cahaya), *Tonik Penenang* (pemulihan Sanity), dan *Penawar Racun*.
- **Kedai Minum:** Tempat beristirahat memulihkan stres warga, mendengar desas-desus pasar, dan merekrut tentara bayaran berpengalaman.

### 7. Pengepungan Gerbang Pemukiman
- **Tingkat Ancaman (Threat Level):** Dihitung otomatis dari akumulasi emas brankas desa, komoditas berharga di gudang, hari berjalan, serta level Balai Desa.
- **Pertahanan Pemukiman:** Ditentukan oleh jumlah Garda Milisi yang bertugas, keberadaan Menara Pengawas, dan Benteng Pertahanan.
- **Dampak Serbuan:**
  - *Kemenangan:* Penjarah dipukul mundur; desa memperoleh rampasan emas dan membebaskan tawanan warga baru.
  - *Kekalahan:* Penjarah menjarah brankas emas, merusak ransum gudang, dan menimbulkan korban jiwa warga.

### 8. Boss Katakombe & Kondisi Kemenangan
Untuk menuntaskan skenario cerita utama dan mengamankan status kemenangan naratif (*Victory Screen*), pemain harus memenuhi **3 syarat utama**:
1. Menemukan dan menumbangkan **Abyssal Overlord** di kedalaman lantai 5 Katakombe.
2. Meningkatkan **Balai Desa ke Level 3**.
3. Membangun infrastruktur pertahanan desa hingga memiliki **Nilai Pertahanan >= 70**.

Setelah memenangkan permainan, pemain dapat memilih untuk melanjutkan permainan dalam **Mode Sandbox Tanpa Batas (*Endless Mode*)**.

### 9. Sistem Simpan & Muat Permainan
- **Penyimpanan Manual:** Tersedia 3 slot penyimpanan terpisah (`slot_1`, `slot_2`, `slot_3`) yang dapat diakses kapan saja dari Markas Pemukiman melalui tombol `[S]`.
- **Penyimpanan Otomatis (Autosave):** Permainan secara otomatis menyimpan data ke slot `autosave` setiap kali pemain melewati hari (`PassDay`) atau menyelesaikan ekspedisi katakombe.

---

## Struktur & Arsitektur Proyek

Proyek ini menggunakan arsitektur Go yang modular dan terstruktur rapi:

```
terminal-odyssey/
├── .github/
│   └── workflows/          # Konfigurasi pipeline GitHub Actions (CI & CD Release)
│       ├── ci.yml          # Verifikasi otomatis (vet, unit test, build test)
│       └── release.yml     # Multi-platform cross-compilation & GitHub Release otomatis
├── cmd/
│   └── odyssey/            # Entry point aplikasi utama (main.go)
├── data/                   # File data JSON dan parser loader bawaan (//go:embed)
│   ├── alchemy/            # Definisi resep dan ramuan alkimia
│   ├── buildings/          # Konfigurasi fasilitas desa dan biaya upgrade
│   ├── companions/         # Roster data pendamping tentara bayaran
│   ├── crafting/           # Resep pembuatan senjata pandai besi
│   ├── dungeon/            # Definisi ruangan, musuh, dan boss Abyssal
│   ├── economy/            # Konfigurasi komoditas dan rute kafilah dagang
│   ├── events/             # Data kejadian acak naratif
│   ├── scenarios/          # Konfigurasi perkemahan awal pemukiman
│   └── siege/              # Kelompok penyerbu dan kalkulasi pengepungan
├── internal/
│   ├── alchemy/            # Engine peracikan ramuan dan inventori obat
│   ├── character/          # Progresi stat karakter, senjata, dan ransel
│   ├── combat/             # Mesin pertempuran taktis berbasis giliran
│   ├── dungeon/            # Generator labirin prosedural, obor, dan kewarasan
│   ├── economy/            # Pasar komoditas, fluktuasi harga, dan kafilah
│   ├── engine/             # Orkestrator pusat logika permainan dan siklus harian
│   ├── save/               # Serialisasi dan persistensi save/load JSON
│   ├── settlement/         # Simulasi pemukiman, warga, pekerja, dan fasilitas
│   ├── siege/              # Mesin kalkulasi ancaman dan resolusi penyerbuan
│   ├── tavern/             # Manajemen kedai minum dan rekrutmen pendamping
│   └── ui/                 # Framework Text User Interface (Bubble Tea & Lip Gloss)
│       ├── styles/         # Palet warna TrueColor dan styling konsisten
│       └── views/          # Modul tampilan individual (lebar standar 78 kolom)
├── Makefile                # Otomasi pengujian, kompilasi, dan cross-compile
├── game-concept.md         # Dokumen rancangan konsep desain permainan
├── improvement-idea.md     # Backlog ide pengembangan jangka panjang
├── go.mod                  # Manifes modul Go
└── README.md               # Dokumentasi utama proyek
```

---

## Otomasi CI/CD & Rilis

Repositori ini telah dilengkapi dengan pipeline otomatisasi GitHub Actions:

1. **Continuous Integration (`.github/workflows/ci.yml`):**
   - Berjalan secara otomatis khusus ketika membuat Pull Request dari branch `dev` menuju branch `main`.
   - Menjalankan pemeriksaan integritas dependensi (`go mod verify`), analisis kode statis (`go vet ./...`), pengujian unit menyeluruh (`go test -v -count=1 ./...`), serta verifikasi kompilasi aplikasi sebelum digabungkan ke `main`.

2. **Continuous Delivery / Release (`.github/workflows/release.yml`):**
   - Terpicu secara otomatis ketika Anda membuat Git Tag rilis versi baru dengan awalan `v*` (contoh: `v1.0.0`).
   - Melakukan *cross-compilation* mandiri (*standalone binary*) untuk berbagai arsitektur dan sistem operasi:
     - Linux (`amd64`, `arm64`)
     - Windows (`amd64` `.exe`)
     - macOS (`Intel amd64`, `Apple Silicon arm64`)
   - Menghasilkan berkas verifikasi `checksums.txt` (SHA256).
   - Mempublikasikan seluruh biner rilis ke halaman resmi **GitHub Releases** beserta catatan rilis otomatis.

---

## Pengujian Kode

Seluruh paket kode dilengkapi dengan pengujian unit komprehensif. Jalankan perintah berikut untuk memverifikasi keandalan sistem:

```bash
# Menjalankan seluruh pengujian unit
make test

# Menjalankan analisis kode statis
make vet

# Menjalankan pengujian dengan metrik cakupan kode
make test-coverage
```

---

## Lisensi

Proyek ini didistribusikan di bawah lisensi terbuka untuk tujuan eksplorasi, edukasi, dan hiburan berbasis perangkat lunak bebas.
