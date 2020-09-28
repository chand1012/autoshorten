from tkinter import TclError, Tk

from tinyurl import shorten

window = Tk() # this is actually hidden
window.withdraw() # this hides the window

shortlink = ''
content = ''
print('Waiting for user to copy URL.')

while True:
    try:
        content = str(window.clipboard_get())
    except TclError as e:
        print(e)
        continue

    if content.startswith('http') and not content is shortlink and len(content) > 30:
        print('User copied link ' + content)
        shortlink = shorten(content)
        print('Made new link ' + shortlink)
        window.clipboard_clear()
        window.clipboard_append(shortlink)
