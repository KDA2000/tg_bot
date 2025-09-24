import re

# Токены и регулярные выражения
token_specification = [
    ('NUMBER', r'\d+(\.\d*)?'),   # число
    ('ID',     r'[a-zA-Z]\w*'),  # идентификатор
    ('PLUS',   r'\+'),            # +
    ('MINUS',  r'-'),             # -
    ('MULT',   r'\*'),            # *
    ('DIV',    r'/'),             # /
    ('LPAREN', r'\('),            # (
    ('RPAREN', r'\)'),            # )
    ('SKIP',   r'[ \t]+'),        # пропуски
    ('MISMATCH', r'.'),           # неверный символ
]

class ParserError(Exception):
    pass

def tokenize(code):
    tok_regex = '|'.join(f'(?P<{name}>{pattern})' for name, pattern in token_specification)
    get_token = re.compile(tok_regex).match
    pos = 0
    mo = get_token(code, pos)
    while mo is not None:
        kind = mo.lastgroup
        value = mo.group(kind)
        if kind == 'NUMBER':
            yield ('number', value)
        elif kind == 'ID':
            yield ('id', value)
        elif kind == 'PLUS':
            yield ('+', value)
        elif kind == 'MINUS':
            yield ('-', value)
        elif kind == 'MULT':
            yield ('*', value)
        elif kind == 'DIV':
            yield ('/', value)
        elif kind == 'LPAREN':
            yield ('(', value)
        elif kind == 'RPAREN':
            yield (')', value)
        elif kind == 'SKIP':
            pass
        elif kind == 'MISMATCH':
            raise RuntimeError(f'Неверный символ {value!r} в позиции {pos}')
        pos = mo.end()
        mo = get_token(code, pos)
    yield ('EOF', '')

class Parser:
    def __init__(self, tokens):
        self.tokens = tokens
        self.next_token()

    def next_token(self):
        self.current_token = next(self.tokens)

    def match(self, expected_type):
        if self.current_token[0] == expected_type:
            val = self.current_token[1]
            self.next_token()
            return val
        else:
            raise ParserError(f"Ошибка! Ожидалось: {expected_type}, получено: {self.current_token[0]}")

    def parse(self):
        tree = self.parseS()
        if self.current_token[0] != 'EOF':
            raise ParserError("Ошибка! Лишние символы в конце входа.")
        return tree

    def parseS(self):
        # S -> T E
        t_node = self.parseT()
        e_node = self.parseEPrime()
        return ('S', t_node, e_node)

    def parseEPrime(self):
        # E -> + T E | - T E | ε
        if self.current_token[0] == '+':
            self.match('+')
            t_node = self.parseT()
            e_node = self.parseEPrime()
            return ('E', '+', t_node, e_node)
        elif self.current_token[0] == '-':
            self.match('-')
            t_node = self.parseT()
            e_node = self.parseEPrime()
            return ('E', '-', t_node, e_node)
        else:
            return ('E', 'ε')

    def parseT(self):
        # T -> F T'
        f_node = self.parseF()
        t_prime_node = self.parseTPrime()
        return ('T', f_node, t_prime_node)

    def parseTPrime(self):
        # T' -> * F T' | / F T' | ε
        if self.current_token[0] == '*':
            self.match('*')
            f_node = self.parseF()
            t_prime_node = self.parseTPrime()
            return ("T'", '*', f_node, t_prime_node)
        elif self.current_token[0] == '/':
            self.match('/')
            f_node = self.parseF()
            t_prime_node = self.parseTPrime()
            return ("T'", '/', f_node, t_prime_node)
        else:
            return ("T'", 'ε')

    def parseF(self):
        # F -> ( S ) | number | id
        if self.current_token[0] == '(':
            self.match('(')
            s_node = self.parseS()
            self.match(')')
            return ('F', '(', s_node, ')')
        elif self.current_token[0] == 'number':
            val = self.match('number')
            return ('F', 'number', val)
        elif self.current_token[0] == 'id':
            val = self.match('id')
            return ('F', 'id', val)
        else:
            raise ParserError("Ошибка! Ожидалось: number, id или '('")

def print_tree(node, indent=0):
    if isinstance(node, tuple):
        print('  ' * indent + str(node[0]))
        for child in node[1:]:
            print_tree(child, indent+1)
    else:
        print('  ' * indent + str(node))

def main():
    expr = input("Введите арифметическое выражение: ")
    try:
        tokens = tokenize(expr)
        parser = Parser(tokens)
        tree = parser.parse()
        print("Выражение корректно.")
        print("Синтаксическое дерево:")
        print_tree(tree)
    except RuntimeError as e:
        print(f"Лексическая ошибка: {e}")
    except ParserError as e:
        print(e)

if __name__ == "__main__":
    main()
