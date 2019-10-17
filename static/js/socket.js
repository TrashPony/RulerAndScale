let ws;

function Connect() {
    ws = new WebSocket("ws://" + window.location.host + "/ws");
    console.log("Websocket - status: " + ws.readyState);

    ws.onopen = function () {
        setInterval(function () {
            ws.send(JSON.stringify({event: "Debug"}));
        }, 1000);
        console.log("Connection opened..." + this.readyState);
    };

    ws.onmessage = function (msg) {
        Draw(JSON.parse(msg.data))
    };

    ws.onerror = function (msg) {
        console.log("Error occured sending..." + msg.data);
    };

    ws.onclose = function (msg) {
        console.log("Disconnected lobby - status " + this.readyState);
    };
}

let Scale = 4;
let left = {};
let right = {};
let back = {};

function Draw(data) {
    console.log(data);

    let canvas = document.getElementById("canvas");
    canvas.width = 400;
    canvas.height = 400;

    let ctx = canvas.getContext("2d");
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // рисуем область измерения
    ctx.strokeStyle = "blue";
    ctx.rect(
        canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)),
        10,
        data.ruler_option.width_max * Scale,
        data.ruler_option.length_max * Scale,
    );
    ctx.stroke();

    // рисуем платформу весов
    ctx.strokeStyle = "black";
    ctx.fillStyle = "rgba(153,161,163,0.71)";
    ctx.roundRect(
        canvas.width / 2 - (data.scale_platform.width * (Scale / 2)),
        10,
        data.scale_platform.width * Scale,
        data.scale_platform.height * Scale,
        5);
    ctx.stroke();
    ctx.fill();

    // рисуем левый датчик
    ctx.fillStyle = "rgb(255,121,0)";
    left = {x: canvas.width / 2 - (data.ruler_option.width_max * (Scale / 2)) - 7, y: 30};
    ctx.fillRect(left.x, left.y, 7, 7);

    // рисуем правый датчик
    right = {x: canvas.width / 2 + (data.ruler_option.width_max * (Scale / 2)), y: 30};
    ctx.fillRect(right.x, right.y, 7, 7);

    // рисуем нижний датчик
    back = {x: canvas.width / 2 - 3.5, y: (data.ruler_option.length_max * (Scale / 2)) * 2 + 10};
    ctx.fillRect(back.x, back.y, 7, 7);


    // рисуем показания
    ctx.strokeStyle = "red";

    ctx.beginPath();       // Начинает новый путь
    ctx.moveTo(left.x + 3.5, left.y + 3.5);    // Рередвигает перо в точку (30, 50)
    ctx.lineTo(left.x + 3.5 + (data.indication.left * Scale), left.y + 3.5);  // Рисует линию до точки (150, 100)
    ctx.stroke();          // Отображает путь

    ctx.beginPath();       // Начинает новый путь
    ctx.moveTo(right.x + 3.5, right.y + 3.5);    // Рередвигает перо в точку (30, 50)
    ctx.lineTo(right.x + 3.5 - (data.indication.right * Scale), left.y + 3.5);  // Рисует линию до точки (150, 100)
    ctx.stroke();          // Отображает путь

    ctx.beginPath();       // Начинает новый путь
    ctx.moveTo(back.x + 3.5, back.y + 3.5);    // Рередвигает перо в точку (30, 50)
    ctx.lineTo(back.x + 3.5, back.y + 3.5 - (data.indication.back * Scale));  // Рисует линию до точки (150, 100)
    ctx.stroke();          // Отображает путь

    //data.indication
}

CanvasRenderingContext2D.prototype.roundRect = function (x, y, w, h, r) {
    if (w < 2 * r) r = w / 2;
    if (h < 2 * r) r = h / 2;
    this.beginPath();
    this.moveTo(x + r, y);
    this.arcTo(x + w, y, x + w, y + h, r);
    this.arcTo(x + w, y + h, x, y + h, r);
    this.arcTo(x, y + h, x, y, r);
    this.arcTo(x, y, x + w, y, r);
    this.closePath();
    return this;
};