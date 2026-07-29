from flask import Flask, request
import redis

app = Flask(__name__)
r = redis.Redis(host='redis', port=6379, db=0)

# Tampilan halaman utama
HTML_VOTE = """
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Aplikasi Voting</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #e2e8f0; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .card { background-color: white; padding: 50px; border-radius: 16px; box-shadow: 0 10px 25px rgba(0,0,0,0.1); text-align: center; max-width: 500px; width: 100%; }
        h1 { color: #2d3748; margin-bottom: 30px; font-size: 28px; }
        .btn-container { display: flex; justify-content: center; gap: 20px; }
        .btn { flex: 1; font-size: 24px; padding: 20px 10px; cursor: pointer; border: none; border-radius: 12px; transition: all 0.3s ease; color: white; font-weight: bold; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .btn-cat { background-color: #f6ad55; }
        .btn-dog { background-color: #63b3ed; }
        .btn:hover { transform: translateY(-5px); box-shadow: 0 8px 15px rgba(0,0,0,0.2); }
        .btn:active { transform: translateY(0); }
    </style>
</head>
<body>
    <div class="card">
        <h1>Pilih Hewan Favoritmu!</h1>
        <form method="POST" class="btn-container">
            <button class="btn btn-cat" name="vote" value="Kucing">🐱 Kucing</button>
            <button class="btn btn-dog" name="vote" value="Anjing">🐶 Anjing</button>
        </form>
    </div>
</body>
</html>
"""

# Tampilan setelah berhasil voting
HTML_SUCCESS = """
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Vote Berhasil</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #e2e8f0; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .card { background-color: white; padding: 50px; border-radius: 16px; box-shadow: 0 10px 25px rgba(0,0,0,0.1); text-align: center; }
        h2 { color: #48bb78; }
        a { display: inline-block; margin-top: 20px; padding: 10px 20px; background-color: #4a5568; color: white; text-decoration: none; border-radius: 8px; transition: 0.3s; }
        a:hover { background-color: #2d3748; }
    </style>
</head>
<body>
    <div class="card">
        <h2>🎉 Suara Anda berhasil dicatat!</h2>
        <a href="/">Kembali Voting</a>
    </div>
</body>
</html>
"""

@app.route("/", methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        vote = request.form['vote']
        r.rpush('votes', vote)
        return HTML_SUCCESS
    return HTML_VOTE

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=80)