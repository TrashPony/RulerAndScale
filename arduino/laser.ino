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
#define RIGHT_PING_LAN 1
#define TOP_PING_LAN   2
#define BACK_PING_LAN  3

// Диоды
#define RED_LED_PIN  9
#define GREEN_LED_PIN  5

#define TOP_MAX          29
#define WIDTH_MAX        20
#define LENGTH_MAX       19

boolean onlyWeight = false;

int widthBox;
int heightBox;
int lengthBox;

boolean debug = false;

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

void setup()
{
  // подключаем лазер
  lox.begin();

  pinMode(RED_LED_PIN, OUTPUT);
  pinMode(GREEN_LED_PIN, OUTPUT);
  pinMode(BUTTON, INPUT);

  Serial.begin(19200);
  // ждем пока откроется сериал порт
  while(!Serial) {}

  // начало работы с I²C хабом
  splitter.begin();

  // ждём одну секунду
  delay(1000);
}

void loop()
{
  if (digitalRead(BUTTON) == LOW) {
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

  Indication();

  if (Serial.available()) {

    byte incomingByte = Serial.read();

    if (incomingByte == 0x95) {
      byte buf[4] = {0x7F, 0x7F, 0x7F, 0x7F};
      Serial.write(buf, 4);
      Serial.flush();
    }

    // запрос габаритов
    if (incomingByte == 0x88) {
      byte buf[12] = {
        0x2D, 0x0B, widthBox, 0x7B,
        0x2D, 0x16, heightBox, 0x7B,
        0x2D, 0x21, lengthBox, 0x7B,
        };
      Serial.write(buf, 12);
    }

    if (incomingByte == 0x66) {
      if (onlyWeight) {
        digitalWrite(GREEN_LED_PIN, HIGH);
        digitalWrite(RED_LED_PIN, LOW);
      } else {
        digitalWrite(RED_LED_PIN, HIGH);
        digitalWrite(GREEN_LED_PIN, LOW);
      }
    }

    if (incomingByte == 0x55) {
      digitalWrite(RED_LED_PIN, LOW);
      digitalWrite(GREEN_LED_PIN, LOW);
    }
  }
}

int getDistance(int pin) {

  // pin - указываем номер сети для лазера откуда брать данные
  splitter.setBusChannel(pin);

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

  if (measure.RangeStatus != 4 && state == 0) {
    return distInt;
  } else {
    return -1;
  }
}

void Indication() {

  int right = getDistance(RIGHT_PING_LAN);
  int left = getDistance(LEFT_PING_LAN);
  int top =  getDistance(TOP_PING_LAN);
  int back = getDistance(BACK_PING_LAN);

  widthBox = WIDTH_MAX - (right + left);
  heightBox = TOP_MAX - top;
  lengthBox = LENGTH_MAX - back;

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
    delay(2500);
  }
}