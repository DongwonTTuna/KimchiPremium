from bs4 import BeautifulSoup
import requests, pyupbit, random, time
from telegram.ext import Updater, CommandHandler
try:
    f = open("/kimchipre.txt",'r')
    f.close()
except:
    f = open("/kimchipre.txt",'w')
    f.close()
updater = Updater('TELEGRAM_CODE',use_context=True)
dispatcher = updater.dispatcher
def get_data():
    url = "https://coinmarketcap.com/ko/currencies/bitcoin/"
    r = requests.get(url).text
    soup = BeautifulSoup(r, 'html.parser')
    elems = str(soup.find_all("div", class_='priceValue___11gHJ'))
    elems = float(elems[elems.find("₩") + 1:elems.find('</')].replace(',', ''))
    price = pyupbit.get_current_price("KRW-BTC")
    return (price - elems) / elems * 100
def add(update, context):
    with open("/kimchipre.txt",'w') as f:
        a = ' '.join(context.args)
        f.write(a)
def help(update):
    updater.bot.sendMessage(USER_ID, "/set 퍼센트 [이상,이하]")
add_h = CommandHandler('set', add)
help_h = CommandHandler('help', help)
dispatcher.add_handler(add_h)
dispatcher.add_handler(help_h)
updater.start_polling()
while True:
    with open("/kimchipre.txt", 'r') as f:
        p = f.read()
    if p == '':
        time.sleep(600)
        continue
    p = p.split(' ')
    a = get_data()
    print(a)
    if p[1] == '이상':
         if a >= float(p[0]):
            for d in range(0, 15):
                updater.bot.sendMessage(USER_ID, "김치 프리미엄이 "+ str(a) +"%입니다!")
            with open("/kimchipre.txt", 'w') as f:
                f.write("")
    elif p[1] == '이하':
        if a <= float(p[0]):
            for d in range(0, 15):
                updater.bot.sendMessage(USER_ID, "김치 프리미엄이 " + str(a) + "%입니다!")
            with open("/kimchipre.txt", 'w') as f:
                f.write("")

    else:
        updater.bot.sendMessage(USER_ID, "잘못 입력하셨습니다.")
    time.sleep(random.randrange(30,60))
