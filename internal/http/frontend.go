package http

var apiKeyRequestTempl = `
        <!DOCTYPE html>
        <html>
            <head>
                <link rel="preconnect" href="https://fonts.googleapis.com">
                <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
                <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
                <style>
                    body {
                        background-color: black;
                        color: white;
                        font-family: 'JetBrains Mono', monospace;
                        padding: 40px;
                        max-width: 800px;
                        margin: 0 auto;
                        line-height: 1.6;
                    }
                    .instructions {
                        background: rgba(255, 255, 255, 0.05);
                        padding: 20px;
                        margin-bottom: 30px;
                        border-left: 3px solid #4CAF50;
                    }
                    .instructions h2 {
                        margin-top: 0;
                        color: #4CAF50;
                        font-size: 1.1rem;
                    }
                    .instructions ol {
                        margin: 0;
                        padding-left: 20px;
                    }
                    .instructions li {
                        margin-bottom: 10px;
                    }
                    .warning {
                        color: #ff4444;
                        font-size: 0.9rem;
                        margin-top: 10px;
                    }
                    .form-group {
                        margin-bottom: 20px;
                    }
                    label {
                        display: block;
                        margin-bottom: 8px;
                        color: #4CAF50;
                    }
                    input[type="email"],
                    input[type="text"] {
                        padding: 12px;
                        width: 100%;
                        max-width: 400px;
                        background: rgba(255, 255, 255, 0.05);
                        border: 1px solid #333;
                        color: white;
                        font-family: 'JetBrains Mono', monospace;
                        font-size: 0.9rem;
                        transition: all 0.3s ease;
                    }
                    input[type="email"]:focus,
                    input[type="text"]:focus {
                        outline: none;
                        border-color: #4CAF50;
                        background: rgba(255, 255, 255, 0.1);
                    }
                    input[type="submit"], .reset-button {
                        background: transparent;
                        color: white;
                        padding: 12px 24px;
                        border: 1px solid #4CAF50;
                        cursor: pointer;
                        font-family: 'JetBrains Mono', monospace;
                        font-size: 0.9rem;
                        text-transform: uppercase;
                        letter-spacing: 1px;
                        position: relative;
                        overflow: hidden;
                        transition: all 0.3s ease;
                    }
                    input[type="submit"]:hover {
                        background: #4CAF50;
                    }
                    .reset-button {
                        border-color: #ff4444;
                        text-decoration: none;
                        display: inline-block;
                    }
                    .reset-button:hover {
                        background: #ff4444;
                    }
                    .error {
                        color: #ff4444;
                        margin-bottom: 15px;
                        padding: 10px;
                        background: rgba(255, 68, 68, 0.1);
                        border-left: 3px solid #ff4444;
                    }
                    .response {
                        color: #4CAF50;
                        margin-bottom: 15px;
                        padding: 10px;
                        background: rgba(76, 175, 80, 0.1);
                        border-left: 3px solid #4CAF50;
                    }
                    .button-group {
                        display: flex;
                        gap: 15px;
                        margin-top: 25px;
                    }
                </style>
            </head>
            <body>
                <div class="instructions">
                    <h2>Quote API key generation</h2>
                    <ol>
                        <li>Enter your email address below to receive a one-time PIN (OTP).</li>
                        <li>Once you receive the OTP, enter it below to generate your API keys.</li>
                        <span class="warning">Note: Generating new API keys will invalidate any existing keys for your account.</span>
                    </ol>
                </div>

                {{if .Error}}
                    <div class="error">{{.Error}}</div>
                {{end}}
                {{if .Response}}
                    <div class="response">{{.Response}}</div>
                {{end}}
                
                <form method="POST" action="/authenticate">
                    <div class="form-group">
                    {{if .Email}}
                        <label for="email">Email:</label>
                        <input readonly type="email" id="email" name="email" value="{{.Email}}" required>
                    {{else}}
                        <label for="email">Email:</label>
                        <input type="email" id="email" name="email" required>
                    {{end}}
                    </div>
                    <div class="form-group">
                        <label for="pin">One Time PIN:</label>
                        <input type="text" id="pin" name="pin">
                    </div>
                    <div class="button-group">
                        <input type="submit" value="Submit">
                        <a href="/authenticate" class="reset-button">Reset</a>
                    </div>
                </form>            
            </body>
        </html>
        `
