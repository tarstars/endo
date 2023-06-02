#!/usr/bin/env python3.10

def main():
    a = input()
    match a:
        case 'aa':
            print('letter aa')
        case _:
            print('some other letter')


if __name__ == '__main__':
    main()
