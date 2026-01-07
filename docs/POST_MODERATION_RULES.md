# Moderação pós-denúncia (Perspective API) — Regras de decisão

Este documento descreve uma política simples de moderação **após denúncia**: o post é publicado sem análise prévia e só é analisado quando um usuário denuncia. A partir dos scores (probabilidades) retornados pelo classificador (ex.: Perspective API), o sistema decide **o que fazer temporariamente** com o conteúdo e **se precisa de revisão humana**.

---

## Objetivo

- **Reduzir risco rapidamente** quando a denúncia chega (principalmente ameaças e ataque a identidade).
- **Evitar remoções incorretas** em casos ambíguos (ex.: ironia, citações, contexto).
- Direcionar conteúdo para **fila de moderação** com prioridade adequada.

---

## Entradas

### Scores do classificador (0.0 a 1.0)

- `THREAT`
- `IDENTITY_ATTACK`
- `SEVERE_TOXICITY`
- `TOXICITY`
- `INSULT`
- `PROFANITY`

Exemplo (resumo):
- `TOXICITY`: 0.825
- `INSULT`: 0.830
- `PROFANITY`: 0.438
- `SEVERE_TOXICITY`: 0.354
- `THREAT`: 0.070
- `IDENTITY_ATTACK`: 0.102

---

## Saídas (Estados do post)

- **VISIBLE**: permanece visível normalmente
- **LIMITED**: permanece visível, mas com restrições (reduz alcance, remove de recomendações, aviso)
- **HIDDEN_PENDING_REVIEW**: oculto (ou visível apenas ao autor) até revisão humana
- **REMOVED**: removido (oculto a qualquer usuário)

---

## Política de decisão (regras)

### 1) Regras duras (risco alto → ocultar imediatamente)

Quando houver risco forte de **ameaça**, **ataque a identidade** ou **toxicidade severa**, o conteúdo deve ser ocultado rapidamente, mesmo antes da revisão humana.

- Se `THREAT >= 0.50` → **HIDDEN_PENDING_REVIEW**
- Se `IDENTITY_ATTACK >= 0.50` → **HIDDEN_PENDING_REVIEW**
- Se `SEVERE_TOXICITY >= 0.70` → **HIDDEN_PENDING_REVIEW**

> Variante mais agressiva (opcional):
> - Se `THREAT >= 0.70` ou `IDENTITY_ATTACK >= 0.70` → **REMOVED** (ou manter oculto + revisão obrigatória)

### 2) Zona cinzenta (risco médio → limitar + revisar)

Se não caiu em regra dura, mas ainda há sinal relevante:

- Se `THREAT` entre `0.35` e `0.50` → **HIDDEN_PENDING_REVIEW**
- Se `IDENTITY_ATTACK` entre `0.35` e `0.50` → **HIDDEN_PENDING_REVIEW**
- Se `SEVERE_TOXICITY` entre `0.45` e `0.70` → **HIDDEN_PENDING_REVIEW**

Caso contrário, avaliar toxicidade geral (score composto).

---

## Score composto (toxicidade geral)

Como `TOXICITY`, `INSULT` e `PROFANITY` se sobrepõem, usamos um score agregado:

```text
composite =
  0.45 * TOXICITY +
  0.35 * INSULT +
  0.20 * PROFANITY