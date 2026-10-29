# career-review

Um processo, em forma de skill para agentes de IA, para revisar a carreira de uma pessoa e reescrever o currículo e o perfil do LinkedIn a partir de fatos com fonte.

[Português](#português) · [English](#english)

## Português

### O que faz

O agente conduz cinco fases:

1. **Apurar os fatos.** Pede onde estão os documentos (currículos antigos, contratos, cartas de promoção, perfil atual do LinkedIn, certificados) e entrevista a pessoa. Monta uma base de fatos em que cada linha tem fonte, e aponta divergências e o que falta confirmar.
2. **Acertar o enquadramento.** Para onde a pessoa quer ir, como se descreve hoje, o que não quer parecer, quais realizações destacar, o que tirar do perfil. Cada resposta vira uma decisão registrada.
3. **Escrever.** Título e "Sobre" do LinkedIn, experiências, projetos, competências e currículo de até 2 páginas, em português e, se quiser, em inglês escrito pelo sentido. Tom profissional e objetivo, sem as marcas típicas de texto gerado por IA.
4. **Revisar.** Uma passada de fidelidade (toda frase tem fonte?) e uma de voz (clichês, jargão, ritmo), com nota. Abaixo de 8, corrige.
5. **Aplicar.** Textos prontos para colar, ou aplicação pelo navegador com autorização, com os cuidados que o LinkedIn exige.

### Como usar

**Qualquer agente com acesso à web:** mande esta mensagem para ele.

```text
Leia https://raw.githubusercontent.com/marcos-toliveira/career-review/main/skills/career-review/SKILL.md e siga o processo para revisar a minha carreira, meu currículo e meu LinkedIn.
```

**Chat sem acesso à web** (ChatGPT, Gemini, Claude.ai): copie o conteúdo de [`skills/career-review/SKILL.md`](skills/career-review/SKILL.md), cole como primeira mensagem e anexe os documentos quando o agente pedir.

**Como skill instalada** (Claude Code, Codex, OpenCode e outros que leem `SKILL.md`):

```sh
git clone https://github.com/marcos-toliveira/career-review.git
mkdir -p ~/.claude/skills && cp -r career-review/skills/career-review ~/.claude/skills/
```

Troque `~/.claude/skills` pela pasta de skills do seu agente.

### Antes de começar

- Separe o que tiver: currículos antigos, carteira de trabalho digital ou contratos, cartas de promoção, certificados e o PDF do seu perfil do LinkedIn.
- Reserve algumas conversas. A primeira fase é a mais longa e a mais importante.
- A sua leitura é a revisão final. Corrija o que não soar como você.
- Os documentos ficam com você e com o seu agente. Este repositório não coleta nada.

## English

### What it does

An agent skill that runs a full career review and rewrites a person's résumé and LinkedIn profile from sourced facts. The agent:

1. **Gathers facts.** Asks where the documents are (old résumés, contracts, promotion letters, the current LinkedIn profile, certificates) and interviews the person. Builds a fact base where every line has a source, and flags conflicts and open questions.
2. **Agrees on positioning.** Target role, how the person describes themselves today, what they don't want to look like, which achievements to feature, what to drop. Every answer is logged as a decision.
3. **Writes.** LinkedIn headline and About, experience entries, projects, skills, and a résumé of up to 2 pages. Portuguese first, with an English version written for meaning rather than translated word by word. Professional, plain tone, without the usual marks of AI-generated text.
4. **Reviews.** One pass for fidelity (does every sentence have a source?) and one for voice (clichés, jargon, rhythm), each scored. Below 8, it revises.
5. **Applies.** Copy-ready text, or browser-based edits with the person's explicit permission and the precautions LinkedIn needs.

The skill itself is written in Brazilian Portuguese. Agents follow it in any language and talk to the person in their own language.

### How to use

Send this to any agent with web access:

```text
Read https://raw.githubusercontent.com/marcos-toliveira/career-review/main/skills/career-review/SKILL.md and follow the process to review my career, my résumé and my LinkedIn profile.
```

Or paste the contents of [`skills/career-review/SKILL.md`](skills/career-review/SKILL.md) into a chat, or install it as a skill as shown in the Portuguese section above.

## Contribuindo / Contributing

Veja / see [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença / License

[MIT](LICENSE)
