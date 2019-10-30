// библиотека для работы с I²C хабом
#include <TroykaI2CHub.h>
// библиотека для управления лазером
#include "Adafruit_VL53L0X.h"

// объект для работы с хабом адрес по умолчанию 0x70
TroykaI2CHub splitter;

// адрес устройства лазера один для всех подсетей
#define LOX2_ADDRESS 0x29

// Кнопки
#define BUTTON 8

// Дальномер
#define LEFT_PING_LAN  0
#define RIGHT_PING_LAN 2
#define TOP_PING_LAN   4
#define BACK_PING_LAN  6

// Диоды
#define RED_LED_PIN  9
#define GREEN_LED_PIN  6

// размеры площади измерений
int TOP_MAX    =      93;
int WIDTH_MAX   =     96;
int LENGTH_MAX   =    55;

boolean onlyWeight = false;
boolean calibrate = true;

int widthBox;
int heightBox;
int lengthBox;

int right;
int left;
int top;
int back;

boolean debug = false;                                                                                          ;

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

void setup()
{

  pinMode(RED_LED_PIN, OUTPUT);
  pinMode(GREEN_LED_PIN, OUTPUT);
  pinMode(BUTTON, INPUT);

  Serial.begin(115200, SERIAL_8E1);
  Serial.setTimeout(100);

  // ждем пока откроется сериал порт
  while(!Serial) {}

  // начало работы с I²C хабом
  splitter.begin();

  delay(1000);

  // открываем соеденение с лазером, но на самом деле нет ¯\_(ツ)_/¯
  // т.к. поднять 4 обьекта для каждого лазера нам не хватает памяти ардуины
  // а калибровать при каждом измерение лазера это будет до 2х секунд на снятие показаний
  // мы поднимаем, одно соеденение для всех лазеров, и пытаемя получать данные через него, но иногда бывает:
  // 1) датчики берут показания предудущего датчика
  // 2) датчик начинает показывать рандомные значения величиной в космос
  // от этого спасает калибровка (поднять соеденение 1 раз для всех датчиков по очереди) и все проблемы пропадают.
  // однако если во время работы например снять датчик и вставить обратно то надо опять калибровать.
  lox.begin();
}

boolean start = false;

void loop() {

  if (start) {
    Indication();
  }

  if (digitalRead(BUTTON) == LOW) {

    if (debug){
        Serial.print("ONLYwEIGHT ");
        Serial.print(" ");
        Serial.println(byte(onlyWeight));
    }

    if (onlyWeight) {
      onlyWeight = false;
      digitalWrite(RED_LED_PIN, HIGH);
      digitalWrite(GREEN_LED_PIN, LOW);
    } else {
      onlyWeight = true;
      digitalWrite(RED_LED_PIN, LOW);
      digitalWrite(GREEN_LED_PIN, HIGH);
    }
    //delay(500);
  }

  if (Serial.available()) {

    byte incomingBytes[2];
    Serial.readBytes(incomingBytes, 2);

    if (incomingBytes[0] == 0x90) {
      TOP_MAX = int(incomingBytes[1]);
      return;
    }

    if (incomingBytes[0] == 0x91) {
      WIDTH_MAX = int(incomingBytes[1]);
      return;
    }

    if (incomingBytes[0] == 0x92) {
      LENGTH_MAX = int(incomingBytes[1]);
      return;
    }

    if (incomingBytes[0] == 0x93) {
      calibrate = true;
      return;
    }

    if (incomingBytes[0] == 0x95) {
      start = true;
      byte buf[1] = {0x7F};
      Serial.write(buf, 1);
      return;
    }

    // запрос габаритов
    if (incomingBytes[0] == 0x88) {

      byte ID = incomingBytes[1];

      byte buf[13] = {
        0x2D, 0x0B, widthBox, 0x7B,
        0x2D, 0x16, heightBox, 0x7B,
        0x2D, 0x21, lengthBox, 0x7B,
        byte(onlyWeight)};

      Serial.write(buf, sizeof(buf));
      return;
    }

    // взятие текущих настроек и показаний линейки
    if (incomingBytes[0] == 0x89) {

      byte buf[41] = {
        0x2D, 0x0B, left, 0x7B,
        0x2D, 0xBB, right, 0x7B,
        0x2D, 0x16, top, 0x7B,
        0x2D, 0x21, back, 0x7B,
        0x2D, 0x0B, WIDTH_MAX, 0x7B,
        0x2D, 0x16, TOP_MAX, 0x7B,
        0x2D, 0x21, LENGTH_MAX, 0x7B,
        0x2D, 0x0B, widthBox, 0x7B,
        0x2D, 0x16, heightBox, 0x7B,
        0x2D, 0x21, lengthBox, 0x7B,
        byte(onlyWeight)};

      Serial.write(buf, sizeof(buf));
      return;
    }

    if (incomingBytes[0] == 0x66) {
      if (onlyWeight) {
        digitalWrite(GREEN_LED_PIN, HIGH);
        digitalWrite(RED_LED_PIN, LOW);
      } else {
        digitalWrite(RED_LED_PIN, HIGH);
        digitalWrite(GREEN_LED_PIN, LOW);
      }
      return;
    }

    if (incomingBytes[0] == 0x55) {
      digitalWrite(RED_LED_PIN, LOW);
      digitalWrite(GREEN_LED_PIN, LOW);
      return;
    }
  }

  while (Serial.available()) {
    Serial.read();
  }
}

int getDistance(int pin) {
  // pin - указываем номер сети для лазера откуда брать данные
  splitter.setBusChannel(pin);

  // если требуется калибровка то поднимаем датчику соеденение заного
  if (calibrate) {
    lox.begin();
  }

  return getIndication();
}

int getIndication() {

  VL53L0X_RangingMeasurementData_t measure;
  lox.rangingTest(&measure, false);

  // проверка на доступность устройства
  Wire.beginTransmission(LOX2_ADDRESS);
  byte state = Wire.endTransmission();

  float dist = round((measure.RangeMilliMeter)/10);
  int distInt = int(dist);

  if (state != 0) {
    return 201; // датчик не подключен
  }

  if (measure.RangeStatus != 4 && distInt <= 200) {
    return distInt;
  } else {
    return 202; // датчик ушел за пределы измерения
  }
}

void Indication() {

  right = getDistance(RIGHT_PING_LAN);
  left = getDistance(LEFT_PING_LAN);
  top =  getDistance(TOP_PING_LAN);
  back = getDistance(BACK_PING_LAN);

  widthBox = WIDTH_MAX - (right + left);
  heightBox = TOP_MAX - top;
  lengthBox = LENGTH_MAX - back;

  calibrate = false;
  if (debug) {
    Serial.print("Right_ping: ");
    Serial.print(right);
    Serial.println("cm");

    Serial.print("Left_ping: ");
    Serial.print(left);
    Serial.println("cm");

    Serial.print("Top_ping: ");
    Serial.print(top);
    Serial.println("cm");

    Serial.print("Back_ping: ");
    Serial.print(back);
    Serial.println("cm");

    Serial.println("///////////////////////////////////////");

    Serial.print("widthBox ");
    Serial.print(" ");
    Serial.println(widthBox);

    Serial.print("heightBox ");
    Serial.print(" ");
    Serial.println(heightBox);

    Serial.print("lengthBox ");
    Serial.print(" ");
    Serial.println(lengthBox);

    Serial.println("///////////////////////////////////////");
    //delay(2500);
  }
}