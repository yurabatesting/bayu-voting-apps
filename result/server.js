const express = require('express');
const { Pool } = require('pg');
const app = express();

// Mengambil kredensial dari Environment Variables
const pool = new Pool({
  user: process.env.POSTGRES_USER,
  host: process.env.DB_HOST,
  password: process.env.POSTGRES_PASSWORD,
  database: process.env.POSTGRES_DB || 'postgres'
});

app.get('/', async (req, res) => {
  try {
    // Mengambil data dari database
    const result = await pool.query('SELECT vote, COUNT(*) as count FROM votes GROUP BY vote');
    
    // Menghitung total suara untuk mencari persentase
    let totalVotes = 0;
    result.rows.forEach(row => {
        totalVotes += parseInt(row.count);
    });

    // Membuat kerangka HTML & CSS (Tampilan Modern)
    let html = `
    <!DOCTYPE html>
    <html lang="id">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Hasil Poling</title>
        <style>
            body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f0f4f8; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
            .card { background-color: white; padding: 40px; border-radius: 20px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); width: 100%; max-width: 550px; }
            h1 { color: #1a202c; text-align: center; margin-bottom: 30px; font-size: 28px; }
            .vote-item { margin-bottom: 25px; }
            .vote-header { display: flex; justify-content: space-between; margin-bottom: 10px; font-size: 18px; font-weight: bold; color: #4a5568; }
            .bar-container { width: 100%; background-color: #e2e8f0; border-radius: 10px; overflow: hidden; height: 28px; box-shadow: inset 0 2px 4px rgba(0,0,0,0.1); }
            .bar { height: 100%; display: flex; align-items: center; padding-left: 10px; color: white; font-weight: bold; font-size: 14px; transition: width 0.5s ease-in-out; }
            .bar-Kucing { background-color: #f6ad55; }
            .bar-Anjing { background-color: #63b3ed; }
            .total { text-align: center; margin-top: 30px; font-size: 16px; color: #718096; border-top: 2px dashed #e2e8f0; padding-top: 20px; }
        </style>
        <script>
            // Halaman akan otomatis refresh setiap 2 detik
            setTimeout(() => location.reload(), 2000);
        </script>
    </head>
    <body>
        <div class="card">
            <h1>📊 Hasil Poling Langsung</h1>
    `;
    
    // Menampilkan pesan jika belum ada vote
    if(result.rows.length === 0) {
        html += `<p style="text-align:center; color:#718096; font-style: italic;">Belum ada suara masuk. Silakan lakukan voting!</p>`;
    } else {
        // Membuat bar persentase untuk setiap pilihan
        result.rows.forEach(row => {
            let count = parseInt(row.count);
            let percentage = totalVotes === 0 ? 0 : Math.round((count / totalVotes) * 100);
            let icon = row.vote === 'Kucing' ? '🐱' : (row.vote === 'Anjing' ? '🐶' : '✨');
            let barClass = `bar-${row.vote}`; // Akan menjadi 'bar-Kucing' atau 'bar-Anjing'
            
            html += `
            <div class="vote-item">
                <div class="vote-header">
                    <span>${icon} ${row.vote}</span>
                    <span>${count} Suara (${percentage}%)</span>
                </div>
                <div class="bar-container">
                    <div class="bar ${barClass}" style="width: ${percentage}%"></div>
                </div>
            </div>`;
        });
        
        html += `<div class="total">Total Seluruh Suara: <b>${totalVotes}</b></div>`;
    }

    html += `
        </div>
    </body>
    </html>
    `;
    
    res.send(html);
  } catch (err) {
    // Tampilan Loading jika Database belum siap
    res.status(500).send(`
        <div style="font-family: sans-serif; text-align: center; margin-top: 20%;">
            <h2 style="color: #718096;">Menyiapkan Database...</h2>
            <p>Silakan tunggu sebentar.</p>
        </div>
        <script>setTimeout(() => location.reload(), 2000);</script>
    `);
  }
});

app.listen(80, () => console.log('Result app berjalan...'));