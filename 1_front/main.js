fetch('https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1')
    .then(response => response.json())
    .then(data => {
        let table = document.createElement('table');
        table.style.width = '100%';
        table.setAttribute('border', '1');
        let header = table.createTHead();
        let row = header.insertRow(0);
        let headers = ["id", "symbol", "name"];
        for(let i = 0; i < headers.length; i++) {
            let cell = row.insertCell(i);
            cell.innerHTML = headers[i];
        }
        let tbody = table.createTBody();
        data.forEach((coin, index) => {
            let row = tbody.insertRow();
            if(index < 5) {
                row.style.background = 'blue';
            }
            if(coin.symbol === 'usdt') {
                row.style.background = 'green';
            }
            for(let i = 0; i < headers.length; i++) {
                let cell = row.insertCell(i);
                cell.innerHTML = coin[headers[i]];
            }
        });
        document.body.appendChild(table);
    });
