Dokumen ini memuat rancangan konsep desain komprehensif untuk game berbasis CLI (Command-Line Interface) modern yang memadukan manajemen pemukiman (settlement management), simulasi ekonomi dan perdagangan (trading simulation), serta penjelajahan bawah tanah naratif prosedural (narrative roguelike dungeon crawler).

# 1. Visi dan Pilar Utama

- **Platform dan Visual:** Berjalan pada Command-Line Interface (CLI) modern berbasis teks. Menggunakan elemen box drawing, simbol Unicode/UTF-8, dan palet warna TrueColor tanpa kontrol pergerakan ubin grid/WASD yang kaku.
- **Teknologi Utama:** Dibangun menggunakan bahasa pemrograman Go (Golang) memanfaatkan framework Text User Interface (TUI) Bubble Tea serta sistem tata letak dan pewarnaan Lip Gloss.
- **Filosofi Gameplay:** Berbasis giliran (turn-based), ramah untuk sesi santai (idle-friendly), berfokus pada kedalaman strategi naratif, dan memiliki nilai putar ulang tinggi (high replayability) melalui generasi konten prosedural.

# 2. Siklus Permainan Utama (Core Loop)

Alur permainan terbagi ke dalam empat fase yang berkesinambungan dan saling memengaruhi stabilitas perkembangan pemain:

1. **Pembangunan dan Persiapan Desa:** Mengatur alokasi tenaga kerja warga, memproses bahan mentah, merawat fasilitas desa, merekrut tentara bayaran, menempa perlengkapan, dan meracik perbekalan ekspedisi.
2. **Ekspedisi Dungeon Naratif:** Menjelajahi reruntuhan bawah tanah berbasis simpul cerita (node-based) bersama rekan pendamping. Menghadapi dilema moral, jebakan tersembunyi, trauma mental, dan pertempuran taktis untuk mengumpulkan artefak serta material tempa langka.
3. **Realisasi Ekonomi dan Pasar:** Menjual hasil jarahan di pasar lokal desa atau mengangkutnya ke kota-kota lain yang mengalami kelangkaan komoditas demi memperoleh margin arbitrase yang tinggi.
4. **Ekspansi dan Riset:** Mengalokasikan keuntungan emas dan material eksotis untuk melatih stat karakter, menaikkan level pandai besi, memperkuat benteng dari serbuan penjarah, dan membuka rute ekspedisi berisiko lebih tinggi.

# 3. Sistem Pembangunan dan Manajemen Desa

Pemain bertindak sebagai pengelola pemukiman yang bermula dari perkemahan darurat sederhana hingga bertransformasi menjadi kota benteng dagang yang tangguh.

## 3.1. Klasifikasi Sumber Daya

| Sumber Daya | Fungsi Utama | Cara Memperoleh |
|---|---|---|
| **Kayu (Lumber)** | Bahan baku pendirian dan perbaikan fasilitas desa serta bahan bakar musim dingin | Dihasilkan secara berkala oleh Penebang Kayu |
| **Batu (Stone)** | Material penguatan dinding benteng dan bangunan tingkat lanjut | Dihasilkan oleh Penambang Batu atau reruntuhan dungeon |
| **Ransum (Food)** | Konsumsi pangan harian warga dan bekal wajib saat ekspedisi | Diproduksi oleh Petani di ladang pemukiman |
| **Koin Emas (Gold)** | Mata uang transaksi pasar, upah pekerja, rekrutmen pendamping, dan latihan stat | Perdagangan komoditas, penjualan relik jarahan, dan pajak |
| **Warga (Settlers)** | Tenaga kerja aktif yang dapat ditugaskan pada berbagai peran produksi dan pertahanan | Bertambah saat Balai Desa ditingkatkan dan reputasi naik |

## 3.2. Penugasan Tenaga Kerja

- **Petani:** Memastikan lumbung ransum tidak mengalami defisit yang dapat memicu kelaparan.
- **Penebang dan Penambang:** Menyuplai material dasar secara otomatis setiap pergantian hari.
- **Pandai Besi:** Mengolah bijih besi menjadi zirah, senjata khusus, perkakas, dan rantai pengunci.
- **Garda Milisi:** Berpatroli menjaga perimeter desa dan memukul mundur pengepungan bandit saat ekspedisi berlangsung.

## 3.3. Struktur Bangunan dan Peningkatan

- **Balai Desa (Town Hall):** Jantung administrasi pemukiman. Peningkatan level memperbesar kapasitas maksimal bangunan dan menarik lebih banyak populasi baru.
- **Gudang Logistik (Storehouse):** Menentukan batas maksimum penampungan komoditas dan hasil tambang agar tidak terbuang percuma.
- **Bengkel Pandai Besi (Blacksmith):** Membuka pembuatan senjata modular, perisai, dan peningkatan zirah pemain.
- **Pusat Latihan (Training Grounds):** Tempat melatih dan meningkatkan stat inti karakter secara permanen menggunakan emas dan ransum.
- **Laboratorium Alkimia (Apothecary):** Tempat mengolah herba liar menjadi salep pemulih HP, penawar racun, serta minyak obor berkualitas tinggi.
- **Kedai Minum (Tavern):** Fasilitas untuk merekrut pendamping petualang, mendengarkan rumor perdagangan wilayah lain, dan memulihkan kondisi mental warga.
- **Pos Kafilah (Caravan Post):** Membuka jaringan ekspedisi logistik otomatis dan gerobak dagang antar-wilayah.

# 4. Atribut Karakter dan Progresi Stat

Karakter pemain memiliki empat atribut inti yang menentukan efektivitas pertempuran dan interaksi naratif di dungeon:

- **Might (Kekuatan Fisik):** Meningkatkan damage serangan fisik senjata melee dan memperluas kapasitas batas slot inventori ransel.
- **Agility (Kelincahan):** Meningkatkan peluang menghindar (dodge), peluang serangan kritikal (crit chance), serta inisiatif giliran pertama dalam duel.
- **Resolve (Ketahanan Mental):** Mengurangi laju penurunan Sanity akibat teror kegelapan, memperbesar batas HP maksimum, dan meningkatkan peluang memicu Kegigihan Heroik.
- **Ingenuity (Kecerdikan):** Membuka opsi dialog cerdas di Ruang Misteri (seperti memecahkan teka-teki kuno), meningkatkan hasil tawar-menawar harga di pasar dagang, dan efisiensi meracik bahan alkimia.

Metode Peningkatan Stat: Melalui pelatihan terstruktur di Pusat Latihan pemukiman, penemuan Buku Pengetahuan Kuno di kedalaman dungeon, atau bonus permanen dari relik artefak langka.

# 5. Sistem Senjata dan Pembuatan Modular (Crafting Engine)

- **Pedang dan Perisai (Sword & Shield):** Gaya bertarung seimbang yang memberikan mekanik Parry dan Block untuk meniadakan kerusakan musuh.
- **Gada dan Palu Berat (Blunt Weapon):** Damage masif dengan kecepatan serangan lambat; memiliki peluang menghasilkan efek Pingsan (Stun) pada musuh.
- **Belati Kembar (Twin Daggers):** Serangan berkecepatan tinggi dengan penskalaan Agility, memicu efek pendarahan (Bleed) berkala.
- **Tombak Panjang (Polearm):** Memberikan keuntungan serangan pendahuluan (Preemptive Strike) sebelum musuh memasuki jarak serang.
- **Parameter Senjata:** Setiap senjata memiliki nilai Attack Damage (ATK), Critical Rate, Inisiatif Kecepatan, dan Durabilitas (ketahanan aus).
- **Sistem Modifikasi Acak (Affix):** Pandai Besi dapat memadukan material monster dungeon untuk memberikan efek elemental khusus (contoh: racun, penghancur zirah, atau peredam teror Sanity).

# 6. Sistem Perdagangan dan Ekonomi Dinamis

- **Dinamika Harga:** Nilai beli dan jual setiap komoditas mengalami fluktuasi harian berkisar antara 5% hingga 35% berdasarkan musim dan stabilitas wilayah.
- **Peluang Arbitrase:** Keuntungan finansial terbesar diperoleh dengan membeli komoditas berlebih di kota produksi berbiaya rendah dan menjualnya ke kota yang mengalami krisis suplai.
- **Kapasitas Logistik:** Daya angkut barang dagang dibatasi oleh kapasitas tas pemain atau jumlah gerobak serta hewan beban yang dimiliki serikat kafilah.
- **Sistem Rumor dan Berita:** Pemain dapat menggali informasi dari percakapan di kedai minum atau pedagang musafir untuk memprediksi perubahan harga sebelum terjadi.

# 7. Sistem Dinamika Lingkungan dan Musim (Seasonal Engine)

| Musim | Efek Pemukiman | Efek Ekspedisi Dungeon | Efek Perdagangan |
|---|---|---|---|
| **Musim Semi (Spring)** | Kelahiran dan migrasi warga meningkat; efisiensi pembibitan ladang naik +20%. | Flora beracun bermunculan di lantai awal; konsumsi ransum normal. | Jalur kafilah terbuka penuh tanpa hambatan cuaca. |
| **Musim Panas (Summer)** | Hasil panen ladang melimpah (+35%), potensi kekeringan air. | Tingkat kehausan meningkat; konsumsi ransum/air naik 1,5x lipat saat menjelajah. | Komoditas pangan mengalami penurunan harga akibat pasokan melimpah. |
| **Musim Gugur (Autumn)** | Musim panen penutup; produksi penebangan kayu meningkat +25%. | Kelembapan tinggi mempercepat redupnya obor sebesar 1,2x lipat. | Harga gandum dan bahan awetan melonjak tinggi di kota-kota lain. |
| **Musim Dingin (Winter)** | Produksi ladang terhenti total; desa memerlukan konsumsi kayu bakar harian. | Suhu beku mengurangi pemulihan di tenda; beberapa lorong tertutup es. | Gerobak dagang antar-kota memakan waktu 2x lebih lama akibat salju. |

# 8. Sistem Rekrutmen Pendamping (Companion & Mercenary)

- **Pencuri (Rogue):** Melucuti perangkap ruangan secara pasif dan membuka peti terkunci tanpa memerlukan Kunci Besi.
- **Sarjana Kuno (Scholar):** Menerjemahkan inskripsi rune purba di Altar Misteri untuk menghindari kutukan serta melipatgandakan temuan relik berharga.
- **Ksatria Pengawal (Vanguard):** Mengalihkan fokus serangan musuh (taunt) dan menyerap sebagian kerusakan fisik tim selama pertempuran.
- **Penyembuh Kelana (Acolyte):** Menurunkan akumulasi stres seluruh anggota tim saat beristirahat dan memulihkan HP tambahan tanpa stok obat.
- **Bagi Hasil Jarahan:** Setiap pendamping menuntut persentase emas (10% hingga 20%) dari total jarahan yang berhasil dibawa pulang.
- **Beban Ransum:** Setiap pendamping menambah kebutuhan ransum ekspedisi sebesar +1 unit ransum per segmen perjalanan.
- **Kematian Permanen:** Pendamping yang gugur di dalam reruntuhan dungeon hilang secara permanen dari daftar rekrutmen dan memicu penurunan moral warga.

# 9. Sistem Dungeon Crawling Naratif

- **Tingkat Nyala Obor:** Dimulai dari 100% dan berkurang 5% hingga 10% setiap kali berpindah ruangan. Obor terang (70-100%) memberi peluang serangan mendadak; remang (30-69%) mengurangi akurasi; gelap gulita (0-29%) memicu ancaman monster pemangsa bayangan.
- **Kapasitas Ransel:** Slot inventori dibatasi (12 hingga 16 slot), menuntut pertimbangan ketat antara membawa perbekalan keselamatan atau membawa pulang relik emas.
- **Sistem Kewarasan (Sanity):** Pengukur mental berskala 0-100. Bila di bawah 30, memicu trauma acak seperti Paranoia, Klaustrofobia, Megalomania, atau Kegigihan Heroik.
- **Generasi Prosedural:** Struktur bawah tanah disusun secara acak dalam Directed Acyclic Graph yang memuat Ruang Pertempuran, Ruang Misteri, Ruang Suaka, dan Titik Evakuasi.
- **Konsekuensi Kekalahan:** Bila HP habis, karakter diselamatkan ke desa dan seluruh ransel ekspedisi hilang, namun kemajuan infrastruktur desa dan kas tetap aman.

# 10. Sistem Pengepungan dan Pertahanan Pemukiman (Siege Events)

- **Tingkat Ancaman (Threat Level):** Dihitung dari akumulasi koin emas di brankas desa, kepenuhan gudang dengan komoditas berharga, serta lamanya durasi ekspedisi yang ditinggalkan tanpa pemimpin.
- **Resolusi Otomatis:** Diselesaikan saat pergantian hari dengan mengadu Daya Serang Musuh melawan Nilai Pertahanan Desa (dinding benteng, menara pengawas, dan jumlah Garda Milisi).
- **Kemenangan vs Kekalahan:** Kemenangan memberi jarahan senjata musuh dan tawanan pekerja baru; kekalahan mengakibatkan pencurian kas emas, kerusakan fasilitas, dan korban jiwa warga.

# 11. Desain Antarmuka Pengguna CLI

## 11.1. Antarmuka Markas Desa

```
+---------------------------------------------------------------------------------+
| DESA: OAKHAVEN [HARI KE-14]                                       Musim: Gugur  |
| Populasi: 8/10 (3 Petani, 2 Penebang, 1 Pandai Besi, 2 Garda)                    |
| Gudang: Kayu: 64/100   Batu: 20/50   Ransum: 18/40                               |
| Kas: 340 Gold                                        Pertahanan: 45              |
+---------------------------------------------------------------------------------+
| LAPORAN HARIAN:                                                                  |
| > Para penebang menghasilkan +15 Kayu (Bonus Musim Gugur).                       |
| > Menara Pengawas mendeteksi jejak bandit di perbatasan utara.                   |
+---------------------------------------------------------------------------------+
| TINDAKAN:                                                                        |
| [1] Masuk ke Pasar dan Perdagangan Komoditas                                     |
| [2] Pembangunan dan Peningkatan Fasilitas Desa                                   |
| [3] Bengkel Pandai Besi (Tempa Senjata & Zirah)                                  |
| [4] Pusat Latihan (Upgrade Stat Karakter)                                        |
| [5] Kunjungi Kedai Minum (Rekrut Pendamping / Cek Rumor)                         |
| [6] Siapkan Perbekalan dan Berangkat ke Katakombe Bawah Tanah                    |
| [SPACE] Lewati Hari (Jalankan siklus produksi harian)                            |
+---------------------------------------------------------------------------------+
```

## 11.2. Antarmuka Bengkel Pandai Besi

```
+---------------------------------------------------------------------------------+
| BENGKEL PANDAI BESI: TEMPA SENJATA                                               |
| Material: Batu: 40 | Biji Besi: 12 | Taring Kalajengking: 2                      |
+---------------------------------------------------------------------------------+
| RESEP TERSEDIA:                                                                  |
| [1] Pedang Baja Standar (Butuh: 4 Biji Besi, 2 Kayu)                             |
|     Stat: 15-20 ATK | +5% Crit                                                   |
|                                                                                   |
| [2] Belati Racun Gurun (Butuh: 2 Besi, 2 Taring)                                 |
|     Stat: 8-12 ATK | +20% Speed | Efek: Poison (3 DMG/turn)                      |
|                                                                                   |
| [3] Perisai Menara Kokoh (Butuh: 6 Besi, 10 Batu)                                |
|     Stat: +15 Armor | Peluang Block 25%                                          |
+---------------------------------------------------------------------------------+
| [Pilih Angka untuk Menempa]   [U] Upgrade Fasilitas Pandai Besi                  |
| [ESC] Kembali ke Pusat Desa                                                      |
+---------------------------------------------------------------------------------+
```

## 11.3. Antarmuka Eksplorasi Bawah Tanah Naratif

```
+---------------------------------------------------------------------------------+
| RUANG 4: ALTAR BERDARAH [Obor: 60%] [Sanity: 45]                                 |
| Pendamping: Ksatria Vanguard (HP 80/100), Pencuri Rogue (Normal)                 |
| Senjata Aktif: Belati Racun Gurun (8-12 ATK, Poison)                             |
|                                                                                   |
| Langkah kakimu membawamu ke aula berkubah bundar. Di tengah                      |
| ruangan, sebuah altar batu hitam memancarkan cahaya ungu                         |
| redup. Tergeletak sebilah belati berukir aksara kuno,                           |
| dikelilingi sisa perlengkapan prajurit yang terbakar.                            |
|                                                                                   |
+---------------------------------------------------------------------------------+
| INVENTORI:                              STATUS EKSPEDISI:                       |
| - Obor (x2)                              HP: 72/100                             |
| - Ransum Kering (x3)                     Kondisi: Waspada                       |
| - Emas Temuan: 85g                       Ransel: 6/12 Slot                      |
+---------------------------------------------------------------------------------+
| TINDAKAN:                                                                        |
| [1] Suruh Pencuri Rogue melucuti kemungkinan jebakan altar                       |
| [2] Ambil belati kuno secara langsung (Uji Ketahanan Mental)                     |
| [3] Abaikan altar dan lanjutkan melangkah menuju Ruang 5                         |
| [4] Putuskan mundur dari ekspedisi dan kembali ke desa                           |
+---------------------------------------------------------------------------------+
```

# 12. Arsitektur Teknis dan Struktur Data (Go)

```go
type GameState int
type Season int
type WeaponType int

const (
    Spring Season = iota
    Summer
    Autumn
    Winter
)

const (
    StateTownMenu GameState = iota
    StateTownBuild
    StateWorkerAssign
    StateBlacksmithCraft
    StateTrainingGrounds
    StateMarketTrade
    StateTavernRecruit
    StateDungeonExplore
    StateCombatTurn
    StateExpeditionSummary
)

type CharacterStats struct {
    Might     int
    Agility   int
    Resolve   int
    Ingenuity int
}

type Weapon struct {
    ID           string
    Name         string
    BaseDamage   [2]int // min, max
    CritRate     float64
    Initiative   int
    Durability   int
    SpecialAffix string
}

type Companion struct {
    Name       string
    Role       string // Vanguard, Rogue, Scholar, Acolyte
    CutPercent int
    HP         int
    MaxHP      int
}

type DungeonNode struct {
    ID          int
    Title       string
    Description string
    Choices     []Choice
    RiskFactor  int
}

type Settlement struct {
    Level      int
    Lumber     int
    Stone      int
    Rations    int
    Treasury   int
    Settlers   int
    DefenseVal int
    Buildings  map[string]int
}

type GameModel struct {
    CurrentState   GameState
    CurrentSeason  Season
    DayCounter     int
    Village        Settlement
    Stats          CharacterStats
    EquippedWeapon Weapon
    Party          []Companion
    ActiveNode     DungeonNode
    PlayerHP       int
    Sanity         int
    TorchMeter     int
    Backpack       []string
}
```
