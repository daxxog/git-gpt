# Git-GPT

Git-GPT is a command-line tool that uses OpenAI's GPT-3.5 model to generate commit messages based on the changes in a git repository.

## Warning: Intellectual Property and Confidential Information Handling

When using this application, please consider the following critical advisory:

This application provides a feature allowing you to send git diffs to the OpenAI API. While this enables powerful collaboration and productivity tools, it has potential privacy and security implications you must be aware of.

Transmitting git diffs to the OpenAI API sends your data to external servers. This data may include confidential or proprietary information. It's crucial to ensure that no sensitive data, such as API keys, passwords, proprietary source code, or any other forms of confidential data, are included in these transmissions.

Please note that as of September 2021, OpenAI's policy was to retain data sent via the API for 30 days. However, this policy may have been updated since then, so it's strongly recommended to review OpenAI's current data usage policy for the most accurate and up-to-date information.

Despite OpenAI's stringent security measures, there are always inherent risks in transmitting data over the internet. By using this feature, you are accepting these risks.

Importantly, this project is licensed under the Apache License 2.0. This means the software is provided "as is," without warranty of any kind, explicit or implied. The license does not assume liability for any damages or losses you may experience, including from the use of this feature.

For corporate users or contributors, please consult with your Information Security department or a legal expert before using this feature. Take the necessary steps to ensure that no intellectual property or confidential information from your company is being inadvertently shared via the API.

Always review and clean your diffs before sharing them through the API, ensuring they do not contain any sensitive data.

Your understanding and compliance are greatly appreciated. Remember to use this feature responsibly, respecting the privacy and security of all sensitive data.


## Installation

Before you can use Git-GPT, you need to make sure you have Go installed on your system. You can download it from the official [Go website](https://golang.org/dl/).

To build the project, navigate to the project's directory and run:

```bash
make build
```

This will generate an executable file `git-gpt` in the project directory.

To install the `git-gpt` command to your system, run:

```bash
make install
```

This will copy the `git-gpt` executable to `/usr/local/bin`, which is typically in the system's PATH, allowing you to run `git-gpt` from anywhere.

## Configuration

Git-GPT needs a valid OpenAI API token to function. You can get this token from your OpenAI account.

Once you have the token, create a configuration file at `~/.config/git-gpt/openai.yaml`:

```yaml
---
token: your-openai-api-token
```

Replace `your-openai-api-token` with your actual OpenAI API token.

## Usage

Git-GPT is used as a drop-in replacement for the `git commit` command.

```bash
git gpt commit -a
```

This will stage all changes, generate a commit message based on those changes, and then open an editor where you can amend the generated commit message.

```bash
git gpt commit -a -m
```

This command does the same as the previous one, but it skips the step where you can amend the generated commit message.

## Makefile Targets

You can use the `make help` command to list all available Makefile targets:

```bash
make help
```

## License

This project is licensed under the Apache License Version 2.0. See the [LICENSE](LICENSE) file for details.
