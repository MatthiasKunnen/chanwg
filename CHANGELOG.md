# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased
See <https://github.com/MatthiasKunnen/chanwg/compare/v2.0.0...master>.

## [2.0.0](https://github.com/MatthiasKunnen/chanwg/compare/v1.1.0...v2.0.0) - 2025-12-09

### <!-- 0 -->Features
- **BREAKING:** Do not close WaitChan before Ready
  (<https://github.com/MatthiasKunnen/chanwg/commit/a9bbe630a1be3c0d4e3cd682b369bd9f60718612>)  
  A call to Ready() is now required to make WaitChan
  complete.

### <!-- 7 -->Miscellaneous Tasks
- Add changelog generator
  (<https://github.com/MatthiasKunnen/chanwg/commit/1e497c88fe7d50406cf4bfc525967cc438820bfe>)

## [1.1.0](https://github.com/MatthiasKunnen/chanwg/compare/v1.0.0...v1.1.0) - 2025-11-20

### <!-- 0 -->Features
- Add WaitGroup.Go
  (<https://github.com/MatthiasKunnen/chanwg/commit/831154b66d84315094ca6a7082706c82af4b587a>)

### <!-- 3 -->Documentation
- Improve examples
  (<https://github.com/MatthiasKunnen/chanwg/commit/180d1f0d517fe7be65c007be6eea016f0b6f6c00>)
- Remove unused IsStarted function in example
  (<https://github.com/MatthiasKunnen/chanwg/commit/06954547b503eeb3fd1336b11d0404325eb073f2>)

### <!-- 6 -->Testing
- Simplify basic test
  (<https://github.com/MatthiasKunnen/chanwg/commit/34de5626a116ddfcadd304abcb0a74074ff7aae3>)
- Add benchmarks and race test
  (<https://github.com/MatthiasKunnen/chanwg/commit/8b6a902bf98b0555f6e9107987a50a7d91ed8ba6>)


## [1.0.0](https://github.com/MatthiasKunnen/chanwg/compare/bd9f4f6ce44d4d0e010431c71fc8627f3ce25f94...v1.0.0) - 2025-07-09

- Implement WaitGroup
