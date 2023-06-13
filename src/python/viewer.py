#!/usr/bin/env python3.10
import collections
import logging

import pygame
import select
import sys

from PIL import Image

logger = logging.Logger(__file__, level=logging.DEBUG)

c_handler = logging.StreamHandler()
c_handler.setLevel(logging.DEBUG)  # Set level for this handler

# Create formatters and add it to handlers
c_format = logging.Formatter('%(name)s - %(levelname)s - %(message)s')
c_handler.setFormatter(c_format)

# Add handlers to the logger
logger.addHandler(c_handler)


def empty_bitmap(shape):
    return pygame.Surface(shape, pygame.SRCALPHA)


class RnaWorkspace:
    def __init__(self, shape):
        logger.info('info')
        logger.debug('debug')
        logger.warning('warning')

        self.black = (0, 0, 0)
        self.red = (255, 0, 0)
        self.green = (0, 255, 0)
        self.blue = (0, 0, 255)
        self.yellow = (255, 255, 0)
        self.magenta = (255, 0, 255)
        self.cyan = (0, 255, 0)
        self.white = (255, 255, 255)
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
        # logger.debug("add color %s", color)
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

    def rotate_ccw(self):
        self.dir = (self.dir[1], -self.dir[0])

    def rotate_cw(self):
        self.dir = (-self.dir[1], self.dir[0])

    def get_color_alpha(self):
        sr = sum(x[0] for x in self.bucket_rgb)
        sg = sum(x[1] for x in self.bucket_rgb)
        sb = sum(x[2] for x in self.bucket_rgb)

        n_color = len(self.bucket_rgb)
        if n_color:
            sr //= n_color
            sg //= n_color
            sb //= n_color

        alpha = sum(self.bucket_alpha)
        n_alpha = len(self.bucket_alpha)

        if n_alpha:
            alpha //= n_alpha
        else:
            alpha = 255

        return sr, sg, sb, alpha

    def line(self):
        color_alpha = self.get_color_alpha()
        try:
            pygame.draw.line(surface=self.bitmaps[-1],
                             color=color_alpha,
                             start_pos=self.pos,
                             end_pos=self.mark,
                             )
        except ValueError:
            print('color argument', color_alpha)

    def get_pixel(self, pos):
        return self.bitmaps[-1].get_at(pos)

    def set_pixel(self, pos, color):
        self.bitmaps[-1].set_at(pos, color)

    def try_fill(self):
        new_color = self.get_color_alpha()
        old_color = self.get_pixel(self.pos)

        logger.debug(f'try fill with {new_color=} {old_color=}')

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
                    if current_color == old_color:
                        yield point
                        visited.add(point)

        q = collections.deque()
        if new_color != old_color:
            q.append(self.pos)

        while q:
            current_point = q.popleft()
            for nei_point in new_nei(*current_point):
                self.set_pixel(nei_point, new_color)
                q.append(nei_point)

    def compose(self):
        logger.debug("compose, number of layers = %d", len(self.bitmaps))
        if len(self.bitmaps) > 1:
            pygame.image.save(self.bitmaps[-1], 'compose_input_0.png')
            pygame.image.save(self.bitmaps[-2], 'compose_input_1.png')
            for y in range(self.shape[1]):
                for x in range(self.shape[0]):
                    p0 = self.bitmaps[-1].get_at((x, y))
                    p1 = self.bitmaps[-2].get_at((x, y))
                    p_res = [0] * 4
                    for z in range(4):
                        p_res[z] = min(255, p0[z] + p1[z] * (255 - p0[3]) // 255)
                    self.bitmaps[-2].set_at((x, y), p_res)
            self.bitmaps.pop()
            pygame.image.save(self.bitmaps[-1], 'compose_output.png')

    def clip(self):
        logger.debug("clip, number of layers = %d", len(self.bitmaps))
        if len(self.bitmaps) > 1:
            for y in range(self.shape[-2]):
                for x in range(self.shape[-1]):
                    p0 = self.bitmaps[-1].get_at((x, y))
                    p1 = self.bitmaps[-2].get_at((x, y))
                    p_res = [0] * 4
                    for z in range(4):
                        p_res[z] = min(255, p1[z] * p0[3] // 255)
                    self.bitmaps[-2].set_at((x, y), p_res)
            self.bitmaps.pop()

    def mark_to_position(self):
        self.mark = self.pos

    def add_new_bitmap(self):
        logger.debug("new layer, number of layers = %d", len(self.bitmaps))
        if len(self.bitmaps) < 10:
            self.bitmaps.append(empty_bitmap(self.shape))

    def process_chunk(self, chunk):
        # logger.debug('chunk %s', chunk)
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
                self.mark_to_position()
            case 'PFFICCP':
                self.line()
            case 'PIIPIIP':
                self.try_fill()
            case 'PCCPFFP':
                self.add_new_bitmap()
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
    chunk_meter = 0
    while running:
        clock.tick(40)
        if select.select([sys.stdin, ], [], [], 0.0)[0]:
            for _ in range(1000):
                chunk = sys.stdin.read(8)
                if chunk:
                    try:
                        ws.process_chunk(chunk[:-1])
                    except ValueError:
                        print(f'{chunk_meter=}')
                        raise
                    chunk_meter += 1
                else:
                    break
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False
        screen.fill((0, 170, 0))

        surface = ws.bitmaps[-1]
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
