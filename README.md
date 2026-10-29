# career-review

Um processo, em forma de skill para agentes de IA, para revisar a carreira de uma pessoa e reescrever o currículo e o perfil do LinkedIn a partir de fatos com fonte.

[Português](#português) · [English](#english)

## Português

### O que faz

O agente conduz cinco fases:

1. **Apurar os fatos.** Pede onde estão os documentos (currículos antigos, contratos, cartas de promoção, perfil atual do LinkedIn, certificados, registros de projetos) e entrevista a pessoa. Monta uma base de fatos em que cada linha tem fonte, mostra as divergências e o que falta confirmar, e só segue depois da validação.
2. **Acertar o enquadramento.** Para onde a pessoa quer ir, como se descreve hoje, como não quer ser percebida, quais realizações destacar, o que tirar do perfil. Cada resposta vira uma decisão registrada.
3. **Escrever.** Título e "Sobre" do LinkedIn, experiências, projetos, competências e currículo, em português e, se quiser, em inglês escrito pelo sentido. Tom profissional e objetivo, sem as marcas típicas de texto gerado por IA.
4. **Revisar.** Uma passada de fidelidade (toda frase tem fonte?) e uma de voz (clichês, jargão, ritmo), com nota. Abaixo de 8, corrige, até duas rodadas.
5. **Entregar.** Textos prontos para a pessoa colar no LinkedIn, com os cuidados de cada campo. O agente não mexe no perfil: o LinkedIn proíbe automação.

Serve para qualquer área. Não depende de código nem de repositórios: as evidências são o que a pessoa entregou, em que escala e com que resultado.

### Como usar

**Agente com acesso à web:** mande esta mensagem.

```text
Leia https://raw.githubusercontent.com/marcos-toliveira/career-review/main/skills/career-review/SKILL.md e siga o processo para revisar a minha carreira, meu currículo e meu LinkedIn.
```

**Chat sem acesso à web** (ChatGPT, Gemini, Claude.ai e outros): copie todo o conteúdo de [`skills/career-review/SKILL.md`](skills/career-review/SKILL.md), cole como primeira mensagem e anexe os documentos quando o agente pedir.

**Como skill instalada** (Claude Code, Codex, OpenCode e outros que leem `SKILL.md`):

```sh
git clone https://github.com/marcos-toliveira/career-review.git
mkdir -p ~/.claude/skills && cp -r career-review/skills/career-review ~/.claude/skills/
```

Troque `~/.claude/skills` pela pasta de skills do seu agente.

### Antes de começar

- Separe o que tiver: currículos antigos, carteira de trabalho digital ou contratos, cartas de promoção, certificados e o PDF do seu perfil do LinkedIn.
- **Privacidade:** antes de anexar, oculte números de documentos, dados bancários, salários e dados de outras pessoas. Não anexe material confidencial de empregador: passe só um resumo. Envie só o necessário. O que você envia fica sujeito às regras do serviço de IA que você usa. Este repositório não coleta nada.
- Reserve algumas conversas. A primeira fase é a mais longa e a mais importante.
- A sua leitura é a revisão final. Corrija o que não soar como você.

## English

### What it does

An agent skill that runs a full career review and rewrites a person's résumé and LinkedIn profile from sourced facts, in five phases:

1. **Gather facts.** Asks where the documents are (old résumés, contracts, promotion letters, the current LinkedIn profile, certificates, project records) and interviews the person. Builds a fact base where every line has a source, shows conflicts and open questions, and moves on only after the person validates it.
2. **Agree on positioning.** Target role, how the person describes themselves today, how they don't want to be perceived, which achievements to feature, what to drop. Every answer is logged as a decision.
3. **Write.** LinkedIn headline and About, experience entries, projects, skills, and a résumé. Portuguese first, with an English version written for meaning rather than translated word by word. Professional, plain tone, without the usual marks of AI-generated text.
4. **Review.** One pass for fidelity (does every sentence have a source?) and one for voice (clichés, jargon, rhythm), each scored. Below 8, it revises, up to two rounds.
5. **Hand over.** Copy-ready text for the person to paste into LinkedIn, with notes for each field. The agent never touches the profile: LinkedIn's User Agreement prohibits automation.

It works for any field. It does not rely on code or repositories: the evidence is what the person delivered, at what scale and with what result.

The skill is written in Brazilian Portuguese. Agents follow it in any language and talk to the person in their own language.

### How to use

**Agent with web access:** send this message.

```text
Read https://raw.githubusercontent.com/marcos-toliveira/career-review/main/skills/career-review/SKILL.md and follow the process to review my career, my résumé and my LinkedIn profile.
```

**Chat without web access** (ChatGPT, Gemini, Claude.ai and others): copy the full contents of [`skills/career-review/SKILL.md`](skills/career-review/SKILL.md), paste it as your first message, and attach documents when the agent asks.

**As an installed skill** (Claude Code, Codex, OpenCode and others that read `SKILL.md`):

```sh
git clone https://github.com/marcos-toliveira/career-review.git
mkdir -p ~/.claude/skills && cp -r career-review/skills/career-review ~/.claude/skills/
```

Replace `~/.claude/skills` with your agent's skills folder.

### Before you begin

- Gather what you have: old résumés, employment records or contracts, promotion letters, certificates, and a PDF of your LinkedIn profile.
- **Privacy:** before attaching files, redact ID numbers, bank details, salaries and other people's data. Don't attach confidential employer material: share a summary instead. Share only what is needed. Whatever you send is subject to the terms of the AI service you use. This repository collects nothing.
- Plan for a few sessions. The first phase is the longest and the most important.
- Your own reading is the final review. Fix anything that doesn't sound like you.

## Contribuindo / Contributing

Veja / see [CONTRIBUTING.md](CONTRIBUTING.md).

## Licença / License

[MIT](LICENSE)
