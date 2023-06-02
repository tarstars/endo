#!/usr/bin/env python3.10
import collections

import pygame
import select
import sys


def empty_bitmap(shape):
    return pygame.Surface(shape, pygame.SRCALPHA)


class RnaWorkspace:
    def __init__(self, shape):
        self.black = (0, 0, 0)
        self.red = (255, 0, 0)
        self.green = (0, 255, 0)
        self.blue = (0, 0, 255)
        self.yellow = (255, 255, 0)
        self.magenta = (255, 0, 255)
        self.cyan = (0, 255, 0)
        self.white = (255, 255, 0)
        self.transparent = 0
        self.opaque = 255
        self.bucket_rgb = []
        self.bucket_alpha = []

        self.pos = (0, 0)
        self.mark = (0, 0)
        self.dir = (1, 0)
        self.bitmaps = []

        self.shape = shape
        self.bitmaps.append(empty_bitmap(shape=shape))

    def add_color(self, color):
        self.bucket_rgb.append(color)

    def add_transparent(self):
        self.bucket_alpha.append(0)

    def add_opaque(self):
        self.bucket_alpha.append(255)

    def drop_bucket(self):
        self.bucket_alpha = []
        self.bucket_rgb = []

    def move(self):
        self.pos = ((self.pos[0] + self.dir[0]) % self.shape[0],
                    (self.pos[1] + self.dir[1]) % self.shape[1]
                    )

    def rotate_cw(self):
        self.dir = (self.dir[1], -self.dir[0])

    def rotate_ccw(self):
        self.dir = (-self.dir[1], self.dir[0])

    def get_color_alpha(self):
        sr = sum(x[0] for x in self.bucket_rgb)
        sg = sum(x[1] for x in self.bucket_rgb)
        sb = sum(x[2] for x in self.bucket_rgb)

        n = len(self.bucket_rgb)
        if n:
            sr //= n
            sg //= n
            sb //= n

        alpha = sum(self.bucket_alpha)
        if n:
            alpha //= n
        else:
            alpha = 255

        return sr, sg, sb, alpha

    def line(self):
        color_alpha = self.get_color_alpha()
        pygame.draw.line(surface=self.bitmaps[0],
                         color=color_alpha,
                         start_pos=self.pos,
                         end_pos=self.mark,
                         )

    def get_pixel(self, pos):
        return self.bitmaps[0].get_at(pos)

    def set_pixel(self, pos, color):
        self.bitmaps[0].set_at(x_y=pos, color=color)

    def try_fill(self):
        new_color = self.get_color_alpha()
        old_color = self.get_pixel(self.pos)

        def nei(x, y):
            yield x - 1, y
            yield x, y - 1
            yield x + 1, y
            yield x, y + 1

        visited = set()

        def new_nei(x, y):
            for point in nei(x, y):
                if 0 <= point[0] < self.shape[0] and 0 <= point[1] < self.shape[1] and point not in visited:
                    current_color = self.get_pixel(point)
                    if current_color != new_color:
                        yield point
                        visited.add(point)

        q = collections.deque()
        if new_color != old_color:
            q.append(self.pos)

        while q:
            current_point = q.popleft()
            for nei_point in new_nei(*current_point):
                self.set_pixel(nei_point, new_color)

    def compose(self):
        if len(self.bitmaps) > 1:
            for y in range(self.shape[1]):
                for x in range(self.shape[0]):
                    p0 = self.bitmaps[0].get_at((x, y))
                    p1 = self.bitmaps[1].get_at((x, y))
                    p_res = [0] * 4
                    for z in range(4):
                        p_res[z] = p0[z] + p1[z] * (255 - p0[3]) // 255
                    self.bitmaps[1].set_at(x_y=(x, y), color=p_res)
            self.bitmaps = self.bitmaps[1:]

    def clip(self):
        if len(self.bitmaps) > 1:
            for y in range(self.shape[1]):
                for x in range(self.shape[0]):
                    p0 = self.bitmaps[0].get_at(x_y=(x, y))
                    p1 = self.bitmaps[1].get_at(x_y=(x, y))
                    p_res = [0] * 4
                    for z in range(4):
                        p_res[z] = p1[z] * p0[-1] // 255
            self.bitmaps = self.bitmaps[1:]

    def process_chunk(self, chunk):
        match chunk:
            case 'PIPIIIC':
                self.add_color(self.black)
            case 'PIPIIIP':
                self.add_color(self.red)
            case 'PIPIICC':
                self.add_color(self.green)
            case 'PIPIICF':
                self.add_color(self.yellow)
            case 'PIPIICP':
                self.add_color(self.blue)
            case 'PIPIIFC':
                self.add_color(self.magenta)
            case 'PIPIIFF':
                self.add_color(self.cyan)
            case 'PIPIIPC':
                self.add_color(self.white)
            case 'PIPIIPF':
                self.add_transparent()
            case 'PIPIIPP':
                self.add_opaque()
            case 'PIIPICP':
                self.drop_bucket()
            case 'PIIIIIP':
                self.move()
            case 'PCCCCCP':
                self.rotate_ccw()
            case 'PFFFFFP':
                self.rotate_cw()
            case 'PCCIFFP':
                self.mark = self.pos
            case 'PFFICCP':
                self.line()
            case 'PIIPIIP':
                self.try_fill()
            case 'PCCPFFP':
                if len(self.bitmaps) < 10:
                    self.bitmaps.append(empty_bitmap(self.shape))
            case 'PFFPCCP':
                self.compose()
            case 'PFFICCF':
                self.clip()
            case _:
                pass


def main():
    shape = (600, 600)
    ws = RnaWorkspace(shape=shape)

    pygame.init()
    screen = pygame.display.set_mode(shape)
    running = True

    clock = pygame.time.Clock()
    while running:
        clock.tick(40)
        if select.select([sys.stdin, ], [], [], 0.0)[0]:
            while True:
                chunk = sys.stdin.read(7)
                if chunk:
                    ws.process_chunk(chunk)
                else:
                    break
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
        screen.fill((0, 170, 0))

        surface = ws.bitmaps[0]
        opaque_surface = surface.convert_alpha()
        for x in range(opaque_surface.get_width()):
            for y in range(opaque_surface.get_height()):
                color = opaque_surface.get_at((x, y))
                opaque_surface.set_at((x, y), (*color[:3], 255))
        screen.blit(opaque_surface, (0, 0))
        pygame.display.update()
    pygame.quit()


if __name__ == "__main__":
    main()