#!/usr/bin/env bash
# Credential-free graph checks and execution of the real aggregate shell only.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
ruby - "${1:-$ROOT/.github/workflows/ci.yml}" <<'RUBY'
require 'yaml'
require 'shellwords'
require 'open3'

NATIVE_TARGETS = [
  %w[ubuntu-24.04 linux amd64], %w[macos-15 darwin arm64],
  %w[macos-15-intel darwin amd64], %w[windows-2025 windows amd64]
].freeze
MAC_TARGETS = NATIVE_TARGETS.select { |target| target[1] == 'darwin' }.freeze
SUITES = %w[
  test-macos-ticket-transport.sh test-apple-fresh-completion.sh
  test-macos-certificate-import.sh test-sign-macos.sh
  test-prepare-macos-signing.sh test-verify-macos-package.sh
].freeze

def assert(condition, message)
  raise message unless condition
end

def expression(value)
  value.to_s.strip.sub(/\A\$\{\{\s*/, '').sub(/\s*\}\}\z/, '').strip
end

# Interpret the runner predicates used by fixture steps, rather than matching
# YAML text. Unknown conditions fail closed instead of assuming a suite runs.
def enabled?(step, target)
  return true unless step.key?('if')
  condition = expression(step['if'])
  return true if condition == 'true'
  return false if condition == 'false'
  match = condition.match(/\A(runner\.os|matrix\.(?:os|goos|goarch))\s*(==|!=)\s*(['"])(.*?)\3\z/)
  assert(match, "unsupported step predicate: #{condition}")
  values = {
    'runner.os' => {'linux' => 'Linux', 'darwin' => 'macOS', 'windows' => 'Windows'}.fetch(target[1]),
    'matrix.os' => target[0], 'matrix.goos' => target[1], 'matrix.goarch' => target[2]
  }
  equal = values.fetch(match[1]) == match[4]
  match[2] == '==' ? equal : !equal
end

def calls?(step, script)
  step.fetch('run', '').lines.any? do |line|
    next false if line.lstrip.start_with?('#')
    words = Shellwords.split(line)
    words.each_cons(2).any? { |command, path| command == 'bash' && path == "scripts/#{script}" }
  rescue ArgumentError
    false # Other shell syntax cannot establish a required suite invocation.
  end
end

def targets(job)
  matrix = job.fetch('strategy').fetch('matrix')
  axes = matrix.reject { |key, _| %w[include exclude].include?(key) }
  rows = axes.empty? ? [] : axes.values.first.product(*axes.values.drop(1)).map { |row| axes.keys.zip(row).to_h }
  rows += matrix.fetch('include', [])
  rows.reject! { |row| matrix.fetch('exclude', []).any? { |excluded| excluded.all? { |key, value| row[key] == value } } }
  rows.map { |row| row.values_at('os', 'goos', 'goarch') }
end

def contract(workflow)
  jobs = workflow.fetch('jobs')
  native = jobs.fetch('native')
  assert(targets(native).sort == NATIVE_TARGETS.sort, 'native target or architecture missing/changed')
  assert(native['name'] == 'Native ${{ matrix.os }}', 'native check names changed')
  fixture_ids = jobs.keys.select do |id|
    id != 'native' && jobs[id].fetch('steps', []).any? { |step| calls?(step, 'test-sign-macos.sh') }
  end
  assert(fixture_ids.length == 1, 'required independent Apple fixture matrix missing/ambiguous')
  fixture_id = fixture_ids.first
  fixtures = jobs.fetch(fixture_id)
  assert(targets(fixtures).sort == MAC_TARGETS.sort, 'both macOS fixture architectures are required')
  assert(Array(fixtures['needs']).empty?, 'Apple fixtures must run independently of native builds')
  [native, fixtures].each do |job|
    assert(!job.key?('if') && !job['continue-on-error'], 'required matrix can be skipped or tolerated')
    assert(job['timeout-minutes'] == 30, 'matrix budget must remain 30 minutes')
    assert(job.fetch('strategy')['fail-fast'] == false, 'matrix must retain every failure witness')
    assert(expression(job['runs-on']) == 'matrix.os', 'runner does not follow the target matrix')
    assert(job.dig('defaults', 'run', 'shell') == 'bash', 'matrix requires Bash semantics')
    assert(!job.key?('environment') && !job.to_s.include?('secrets.'), 'credential boundary changed')
    assert(job.fetch('permissions', workflow['permissions']) == {'contents' => 'read'}, 'matrix permissions changed')
    checkout = job.fetch('steps').find { |step| step.fetch('uses', '').start_with?('actions/checkout@') }
    assert(checkout && checkout['uses'] == 'actions/checkout@d23441a48e516b6c34aea4fa41551a30e30af803', 'checkout pin changed')
    assert(checkout.dig('with', 'persist-credentials') == false && checkout.dig('with', 'fetch-depth') == 2, 'checkout custody changed')
    assert(expression(checkout.dig('with', 'ref')) == 'inputs.checkout_ref || github.sha', 'checkout source binding changed')
    setup = job.fetch('steps').find { |step| step.fetch('uses', '').start_with?('actions/setup-go@') }
    assert(setup && setup['uses'] == 'actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16', 'Go action pin changed')
    assert(setup.dig('with', 'go-version') == '1.26.10', 'Go toolchain changed')
  end
  assert(fixtures['env'] == native['env'], 'fixture environment differs from native environment')
  assert(fixtures['steps'].none? { |step| step.key?('uses') && !step['uses'].start_with?('actions/checkout@', 'actions/setup-go@') }, 'fixtures must synthesize their own inputs')
  assert(fixtures['steps'].none? { |step| step.fetch('run', '').match?(/ci-native\.sh|package-macos\.sh|agent-heavy/) }, 'fixtures contain native builds or a local queue wrapper')
  MAC_TARGETS.each do |target|
    SUITES.each do |suite|
      steps = fixtures['steps'].select { |step| calls?(step, suite) && enabled?(step, target) }
      assert(steps.length == 1 && !steps.first['continue-on-error'], "#{target[0]}: missing, duplicate or tolerated #{suite}")
      body = steps.first.fetch('run')
      assert(body.include?('::group::') && body.include?('::endgroup::') && body.include?('TIMEFORMAT=') && body.match?(/^\s*time\s/), "#{suite}: missing group/timing labels")
      assert(native['steps'].none? { |step| calls?(step, suite) && enabled?(step, target) }, "#{suite}: expensive macOS fixture co-located with native build")
    end
  end
  SUITES.first(2).each do |suite|
    assert(native['steps'].any? { |step| calls?(step, suite) && enabled?(step, NATIVE_TARGETS.first) }, "Linux regression omitted: #{suite}")
  end
  assert(native['steps'].any? { |step| calls?(step, 'test-ci-topology.sh') && enabled?(step, NATIVE_TARGETS.first) }, 'Linux does not validate CI topology')

  gate = jobs.fetch('ci-required')
  assert(gate['name'] == 'ci-required' && expression(gate['if']) == 'always()', 'required aggregate name/always condition changed')
  assert(Array(gate['needs']).sort == ['native', fixture_id].sort, 'aggregate must depend on native AND Apple fixtures')
  assert(!gate['continue-on-error'], 'aggregate failure tolerated')
  steps = gate.fetch('steps')
  assert(!steps.empty? && steps.all? { |step| step['run'] && !step.key?('if') && !step['continue-on-error'] }, 'aggregate shell can be skipped')
  bindings = gate.fetch('env').transform_values { |value| expression(value) }
  %w[native].push(fixture_id).each do |id|
    assert(bindings.value?("needs.#{id}.result"), "aggregate does not bind #{id} result")
  end
  # Execute the actual gate body, including the real result-variable bindings.
  # nil removes the variable; empty is distinct. Only success/success may pass.
  states = ['success', 'failure', 'cancelled', 'skipped', '', nil]
  states.product(states).each do |native_result, fixture_result|
    results = {'native' => native_result, fixture_id => fixture_result}
    env = bindings.transform_values do |binding|
      match = binding.match(/\Aneeds\.([\w-]+)\.result\z/)
      assert(match, "unsupported aggregate env binding: #{binding}")
      results.fetch(match[1])
    end
    passed = steps.all? do |step|
      _out, _err, status = Open3.capture3(env, 'bash', '-e', '-o', 'pipefail', '-c', step.fetch('run'))
      status.success?
    end
    expected = native_result == 'success' && fixture_result == 'success'
    assert(passed == expected, "aggregate misclassified #{[native_result, fixture_result].inspect}")
  end
  fixture_id
end

workflow = YAML.load_file(ARGV.fetch(0))
begin
  fixture_id = contract(workflow)
rescue StandardError => error
  warn "FAIL: #{error.message}"
  exit 1
end
puts 'PASS: four native targets, independent dual-Mac fixture coverage, source/security/budget contracts'
puts 'PASS: actual ci-required shell rejects all 35 non-success pairs (including empty/unset)'

# Prove that the checker detects graph/coverage regressions, not just this file's
# spelling. Each mutant changes one behavior in a parsed copy of the workflow.
mutants = {
  'omitted Intel native' => ->(w) { w['jobs']['native']['strategy']['matrix']['include'].reject! { |row| row['os'] == 'macos-15-intel' } },
  'omitted ARM native' => ->(w) { w['jobs']['native']['strategy']['matrix']['include'].reject! { |row| row['os'] == 'macos-15' } },
  'omitted Intel fixtures' => ->(w) { w['jobs'][fixture_id]['strategy']['matrix']['include'].pop },
  'serialized fixtures' => ->(w) { w['jobs'][fixture_id]['needs'] = 'native' },
  'co-located macOS fixtures' => ->(w) { w['jobs']['native']['steps'] << w['jobs'][fixture_id]['steps'].find { |s| calls?(s, SUITES.last) } },
  'missing aggregate dependency' => ->(w) { w['jobs']['ci-required']['needs'] = 'native' },
  'aggregate bypass' => ->(w) { w['jobs']['ci-required']['steps'] = [{'run' => 'exit 0'}] },
  'aggregate OR instead of AND' => ->(w) { w['jobs']['ci-required']['steps'] = [{'run' => '[[ $NATIVE_RESULT == success || $APPLE_FIXTURES_RESULT == success ]]'}] },
  'aggregate empty accepted' => ->(w) { w['jobs']['ci-required']['steps'] = [{'run' => '[[ ${NATIVE_RESULT:-success} == success && ${APPLE_FIXTURES_RESULT:-success} == success ]]'}] }
}
SUITES.each do |suite|
  mutants["omitted #{suite}"] = ->(w) { w['jobs'][fixture_id]['steps'].reject! { |s| calls?(s, suite) } }
  mutants["skipped #{suite}"] = ->(w) { w['jobs'][fixture_id]['steps'].find { |s| calls?(s, suite) }['if'] = "runner.os == 'Linux'" }
end
mutants.each do |label, mutate|
  copy = Marshal.load(Marshal.dump(workflow))
  mutate.call(copy)
  rejected = false
  begin
    contract(copy)
  rescue StandardError
    rejected = true
  end
  assert(rejected, "checker accepted mutant: #{label}")
  puts "PASS: rejected #{label}"
end
RUBY
