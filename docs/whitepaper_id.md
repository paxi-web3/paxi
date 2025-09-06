# Whitepaper Paxi (Bahasa Indonesia)

> **Catatan:** Ini adalah terjemahan komunitas dari *Paxi Whitepaper*. Dokumen resmi tetap berbahasa Inggris. Versi ini dibuat untuk memudahkan pembaca Indonesia memahami visi dan teknologi Paxi.

---

## Pendahuluan

Di dunia di mana teknologi blockchain sering terasa sulit diakses — dipenuhi hambatan teknis dan didominasi oleh whale serta institusi besar — Paxi hadir dengan mimpi yang lebih besar. Bayangkan sebuah ekosistem blockchain yang cepat, aman, dan benar-benar terdesentralisasi. Sebuah protokol yang dirancang bukan hanya untuk segelintir elit, tetapi untuk semua orang.  

Dibangun dengan **Cosmos SDK** dan **mesin konsensus CometBFT**, serta diimplementasikan dalam bahasa pemrograman Go, Paxi menghadirkan infrastruktur blockchain yang efisien, berbiaya rendah, dan ramah pengguna. Paxi juga mendukung **CosmWasm**, sebuah platform smart contract berbasis WebAssembly (Wasm), sehingga pengembang dapat membangun kontrak aman dan efisien menggunakan Rust — ideal untuk DeFi, DAO, dan dApp lintas rantai.  

Filosofi Paxi sederhana: *less is more*. Setiap fitur hanya ada jika benar-benar memberi manfaat dan efisiensi. Fokus pada hal esensial menjadikan Paxi pengalaman yang mulus bagi pengembang, validator, maupun pengguna akhir.  

Misi kami jelas: memberdayakan individu, komunitas, dan bisnis dengan jaringan yang benar-benar terdesentralisasi, adil, dan inklusif.  

---

## Visi

Paxi membayangkan ekosistem blockchain yang **cepat, aman, dan sepenuhnya terdesentralisasi** — bukan milik segelintir whale atau institusi besar, melainkan **milik semua orang**: pengguna, pengembang, dan partisipan.  

Melalui staking, setiap orang bisa ikut konsensus. Melalui DAO, setiap orang bisa ikut mengambil keputusan. Melalui pengembangan aplikasi, setiap orang bisa berinovasi.  

---

## Prinsip Inti

- **Kesederhanaan**: Desain minimalis, API bersih, tanpa logika protokol berlebihan  
- **Kecepatan**: Kinerja tinggi, ribuan transaksi per detik  
- **Desentralisasi**: Validator set terbuka, ambang partisipasi rendah, tata kelola oleh banyak orang  
- **Keamanan**: Berbasis Byzantine Fault Tolerance, kode ramping dan mudah diaudit  
- **Aksesibilitas**: Ramah pengembang dari semua level, termasuk no-code/low-code  

---

## Tumpukan Teknologi

- **Cosmos SDK**: Framework modular untuk blockchain  
- **CometBFT**: Konsensus BFT dengan finalitas cepat  
- **CosmWasm**: Smart contract berbasis Wasm (Rust)  
- **Bahasa Go**: Sederhana, cepat, mendukung konkruensi  

---

## Partisipasi Validator

Paxi menurunkan hambatan teknis & ekonomi agar siapa pun bisa jadi validator.  
- Hanya butuh **1.000 PAXI** untuk menjadi validator  
- 50% validator dipilih berdasarkan voting power tertinggi  
- 50% lainnya dipilih secara **acak berbobot stake**  

Model ini memastikan validator besar tetap ada, tapi validator kecil juga punya kesempatan adil → desentralisasi meningkat.  

---

## Mekanisme Inflasi

- Pasokan awal: **100 juta PAXI**  
- Inflasi:  
  - Tahun 1: maks 8%  
  - Tahun 2: tetap 4%  
  - Tahun 3+: dibatasi 2%  

Distribusi reward blok:  
- 95% ke validator + delegator  
- 5% ke komunitas DAO  

---

## Distribusi Token

| Kategori                    | Alokasi | Jadwal Rilis | Tujuan |
|-----------------------------|---------|--------------|--------|
| Tim & Penasehat             | 15%     | Vesting 2 tahun | Komitmen jangka panjang |
| Yayasan Paxi                | 10%     | Vesting 2 tahun | Pemeliharaan, hukum, branding |
| DAO Paxi                    | 5%      | Unlock awal | Dana ekosistem & tata kelola |
| Investor Privat & Strategis | 15%     | Vesting 2 tahun | Kemitraan & listing |
| ICO                         | 35%     | Unlock fase | Likuiditas & partisipasi publik |
| Insentif & Promosi          | 20%     | Dinamis | Pertumbuhan komunitas |

---

## DAO Paxi

DAO memungkinkan tata kelola on-chain:  
- Voting parameter (inflasi, fee, staking)  
- Upgrade perangkat lunak  
- Pendanaan proyek komunitas  
- Manajemen izin smart contract  

---

## Kasus Penggunaan

- **DeFi**: Infrastruktur murah & cepat  
- **GameFi**: Mesin NFT & gaming  
- **Identitas Sosial**: Sistem kepercayaan tanpa perantara  
- **Enterprise/IoT**: Cocok untuk perangkat ringan  

---

## AMM On-chain

- Modul AMM native (mirip Uniswap V2)  
- Swap PAXI <-> PRC20 dengan biaya rendah (default 0,4%)  
- Pool likuiditas permissionless  
- Reward ke LP dapat diklaim kapan saja  

---

## Pengalaman Pengembang

- SDK & API intuitif  
- IDE khusus Paxi  
- Opsi no-code/low-code  
- Dokumentasi & tutorial lengkap  

---

## Kesimpulan

Paxi bukan sekadar blockchain lain. Paxi adalah gerakan menuju jaringan yang sederhana, cepat, adil, dan terbuka untuk semua.  

**Bangun lebih sedikit. Capai lebih banyak. Bangun di Paxi.**
