extern "C" {
    #include "jieba.h"
}

#include "cppjieba/Jieba.hpp"

#include <iostream>
#include <sstream>
#include <string>

static char** ConvertWords(const std::vector<std::string>& words) {
  char ** res = (char**)malloc(sizeof(char*) * (words.size() + 1));
  for (size_t i = 0; i < words.size(); i++) {
    res[i] = (char*)malloc(sizeof(char) * (words[i].length() + 1));
    strcpy(res[i], words[i].c_str());
  }
  res[words.size()] = NULL;
  return res;
}

static Word* ConvertWords(const std::vector<cppjieba::Word>& words) {
  Word* res = (Word*)malloc(sizeof(Word) * (words.size() + 1));
  for (size_t i = 0; i < words.size(); i++) {
    res[i].offset = words[i].offset;
    res[i].len = words[i].word.size();
  }
  res[words.size()].offset = 0;
  res[words.size()].len = 0;
  return res;
}

static struct CWordWeight* ConvertWords(const std::vector<std::pair<std::string, double> >& words) {
  struct CWordWeight* res = (struct CWordWeight*)malloc(sizeof(struct CWordWeight) * (words.size() + 1));
  for (size_t i = 0; i < words.size(); i++) {
    res[i].word = (char*)malloc(sizeof(char) * (words[i].first.length() + 1));
    strcpy(res[i].word, words[i].first.c_str());
    res[i].weight = words[i].second;
  }
  res[words.size()].word = NULL;
  return res;
}

Jieba NewJieba(const char* dict_path,
      const char* hmm_path, 
      const char* user_dict,
      const char* idf_path,
      const char* stop_words_path) {
  return (Jieba)(new cppjieba::Jieba(dict_path, hmm_path, user_dict, idf_path, stop_words_path));
}

Jieba NewJiebaWithContent(const char* dict_content,
      const char* hmm_content,
      const char* user_content,
      const char* idf_content,
      const char* stop_words_content) {
      std::stringstream dict_ss(dict_content);
      std::stringstream hmm_ss(hmm_content);
      std::stringstream user_ss(user_content);
      std::stringstream idf_ss(idf_content);
      std::stringstream stop_words_ss(stop_words_content);
  return (Jieba)(new cppjieba::Jieba(dict_ss, hmm_ss, user_ss, idf_ss, stop_words_ss));
}

void FreeJieba(Jieba x) {
  delete (cppjieba::Jieba*)x;
}

char** Cut(Jieba x, const char* sentence, int is_hmm_used) {
  std::vector<std::string> words;
  ((cppjieba::Jieba*)x)->Cut(sentence, words, is_hmm_used);
  char** res = ConvertWords(words);
  return res;
}

char** CutAll(Jieba x, const char* sentence) {
  std::vector<std::string> words;
  ((cppjieba::Jieba*)x)->CutAll(sentence, words);
  char** res = ConvertWords(words);
  return res;
}

char** CutForSearch(Jieba x, const char* sentence, int is_hmm_used) {
  std::vector<std::string> words;
  ((cppjieba::Jieba*)x)->CutForSearch(sentence, words, is_hmm_used);
  char** res = ConvertWords(words);
  return res;
}

char** Tag(Jieba x, const char* sentence) {
  std::vector<std::pair<std::string, std::string> > result;
  ((cppjieba::Jieba*)x)->Tag(sentence, result);
  std::vector<std::string> words;
  words.reserve(result.size());
  for (size_t i = 0; i < result.size(); ++i) {
    words.push_back(result[i].first + "/" + result[i].second);
  }
  return ConvertWords(words);
}

void AddWord(Jieba x, const char* word) {
  ((cppjieba::Jieba*)x)->InsertUserWord(word);
}

void AddWordEx(Jieba x, const char* word, int freq, const char* tag) {
  ((cppjieba::Jieba*)x)->InsertUserWord(word, freq, tag);
}

void RemoveWord(Jieba x, const char* word) {
  ((cppjieba::Jieba*)x)->DeleteUserWord(word);
}

Word* Tokenize(Jieba x, const char* sentence, TokenizeMode mode, int is_hmm_used) {
  std::vector<cppjieba::Word> words;
  switch (mode) {
    case SearchMode:
      ((cppjieba::Jieba*)x)->CutForSearch(sentence, words, is_hmm_used);
      return ConvertWords(words);
    default:
      ((cppjieba::Jieba*)x)->Cut(sentence, words, is_hmm_used);
      return ConvertWords(words);
  }
}

struct CWordWeight* ExtractWithWeight(Jieba handle, const char* sentence, int top_k) {
  std::vector<std::pair<std::string, double> > words;
  ((cppjieba::Jieba*)handle)->extractor.Extract(sentence, words, top_k);
  struct CWordWeight* res = ConvertWords(words);
  return res;
}

void FreeWordWeights(struct CWordWeight* wws) {
  struct CWordWeight* x = wws;
  while (x && x->word) {
    free(x->word);
    x->word = NULL;
    x++;
  }
  free(wws);
}

char** Extract(Jieba handle, const char* sentence, int top_k) {
  std::vector<std::string> words;
  ((cppjieba::Jieba*)handle)->extractor.Extract(sentence, words, top_k);
  char** res = ConvertWords(words);
  return res;
}

void printStr(const char *str) {
    std::cout<< "we get :" << str << std::endl;
}