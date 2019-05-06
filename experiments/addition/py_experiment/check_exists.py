"""
Use projected gradient descent to uncover factorized
quantum circuits for addition.
"""

import numpy as np
import torch
import torch.nn as nn
import torch.optim as optim

SUM_BITS = 5
NUM_BITS = SUM_BITS * 2


def main():
    target_matrix = compute_target_matrix()
    forward = ComplexMatrix.random(16)
    middle = ComplexMatrix.random(16)
    backward = ComplexMatrix.random(16)
    forward2 = ComplexMatrix.random(16)
    middle2 = ComplexMatrix.random(16)
    backward2 = ComplexMatrix.random(16)

    # Check that expanding still produces a unitary matrix.
    #     x = forward.expander(5, [0, 1, 2, 3])()
    #     print(x.mul(x.H()).real)

    matrices = [forward, middle, backward, forward2, middle2, backward2]
    expanders = [
        sliding_expander(forward),
        middle.expander(NUM_BITS, list(range(NUM_BITS - 4, NUM_BITS))),
        sliding_expander(backward, forward=False),

        sliding_expander(forward2),
        middle2.expander(NUM_BITS, list(range(NUM_BITS - 4, NUM_BITS))),
        sliding_expander(backward2, forward=False),
    ]
    sgd = optim.SGD([(m.real, m.imag)[i] for m in matrices for i in [0, 1]], lr=200)

    while True:
        product = expanders[0]()
        for e in expanders[1:]:
            product = e().mul(product)
        all_diffs = torch.pow(target_matrix - product.real, 2)
        approx_loss = torch.mean(all_diffs * torch.rand_like(all_diffs))
        exact_loss = torch.mean(all_diffs)
        sgd.zero_grad()
        approx_loss.backward()
        sgd.step()
        print('loss=%.8f' % exact_loss.item())
        for m in matrices:
            m.orthogonalize()


class ComplexMatrix:
    def __init__(self, real, imag):
        self.real = real
        self.imag = imag

    @classmethod
    def random(cls, size):
        res = cls(nn.Parameter(torch.randn(size, size)),
                  nn.Parameter(torch.randn(size, size)))
        res.orthogonalize()
        return res

    @classmethod
    def eye(cls, size):
        eye = torch.eye(size)
        return cls(eye, torch.zeros_like(eye))

    def H(self):
        return ComplexMatrix(self.real.transpose(1, 0), -self.imag.transpose(1, 0))

    def expander(self, num_bits, bit_indices):
        """
        Generate a function that expands this matrix.
        """
        exp = expand_operator(num_bits, bit_indices)

        def expand():
            return ComplexMatrix(exp(self.real), exp(self.imag))

        return expand

    def mul(self, other):
        """
        Generate (self @ other).
        """
        return ComplexMatrix(torch.matmul(self.real, other.real) -
                             torch.matmul(self.imag, other.imag),
                             torch.matmul(self.real, other.imag) +
                             torch.matmul(self.imag, other.real))

    def orthogonalize(self):
        """
        Project this matrix onto the space of unitary
        matrices.
        """
        real = self.real.detach().cpu().numpy()
        imag = self.imag.detach().cpu().numpy()
        full = real.astype(np.complex) + 1j * imag.astype(np.complex)
        u, _, vh = np.linalg.svd(full)
        orthog = np.dot(u, vh)
        real = np.real(orthog).astype(np.float32)
        imag = np.imag(orthog).astype(np.float32)
        self.real.data = torch.from_numpy(real)
        self.imag.data = torch.from_numpy(imag)


def sliding_expander(matrix, forward=True):
    expanders = []
    for i in (range(0, NUM_BITS - 3, 2) if forward else range(NUM_BITS - 4, -1, -2)):
        expanders.append(matrix.expander(NUM_BITS, [i, i+1, i+2, i+3]))

    def fn():
        res = ComplexMatrix.eye(1 << NUM_BITS)
        for x in expanders:
            res = x().mul(res)
        return res

    return fn


def compute_target_matrix():
    """
    Compute the ground-truth unitary matrix for addition.
    """
    target_np = np.zeros([1 << NUM_BITS, 1 << NUM_BITS], dtype=np.float32)
    for i in range(1 << NUM_BITS):
        n1 = 0
        n2 = 0
        for j in range(SUM_BITS):
            n1 |= (0 if i & (1 << (2 * j)) == 0 else 1) << j
            n2 |= (0 if i & (1 << (2 * j + 1)) == 0 else 1) << j
        result = (n1 + n2) & ((1 << SUM_BITS) - 1)
        target_idx = 0
        for j in range(SUM_BITS):
            if n1 & (1 << j) != 0:
                target_idx |= (1 << (2 * j))
            if result & (1 << j) != 0:
                target_idx |= (1 << (2 * j + 1))
        target_np[target_idx, i] = 1
    return torch.from_numpy(target_np)


def expand_operator(num_bits, bit_indices):
    """
    Generate a function that takes small unitary matrices
    to larger unitary matrices that act on the bits in the
    bit_indices list.

    Args:
        num_bits: the number of bits in the bigger matrix.
        bit_indices: indices mapping bits from the smaller
            matrix to those of the bigger one.

    Returns:
        A function that takes a small matrix and produces
        a larger one.
    """
    bit_mask = 0
    for idx in bit_indices:
        bit_mask |= 1 << idx
    inv_mask = ((1 << num_bits) - 1) ^ bit_mask

    def small_idx(large_idx):
        num = 0
        for i, idx in enumerate(bit_indices):
            if large_idx & (1 << idx) != 0:
                num |= 1 << i
        return num

    source_inds = []
    zero_index = 1 << (2 * len(bit_indices))
    for row in range(1 << num_bits):
        for col in range(1 << num_bits):
            if row & inv_mask != col & inv_mask:
                source_inds.append(zero_index)
            else:
                source_idx = small_idx(col)
                dest_idx = small_idx(row)
                source_inds.append(source_idx + dest_idx * (1 << len(bit_indices)))

    def fn(small_matrix):
        flat = torch.cat([small_matrix.view(-1), torch.zeros_like(small_matrix[0, 0:1])])
        flat_big = flat[source_inds]
        return flat_big.view(1 << num_bits, 1 << num_bits)

    return fn


if __name__ == '__main__':
    main()
